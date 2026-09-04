package featureflags

import (
	"fmt"
	"testing"
)

func TestFeatureService_GlobalOff(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: false},
	})

	if fs.IsEnabled("online_booking", "user-123") {
		t.Error("Feature harus OFF ketika Enabled = false")
	}
}

func TestFeatureService_ZeroPercentRollout(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 0},
	})

	if fs.IsEnabled("online_booking", "user-123") {
		t.Error("Feature harus OFF ketika RolloutPercentage = 0")
	}
}

func TestFeatureService_InternalRollout(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, InternalUsers: []string{"internal-team"}},
	})

	if !fs.IsEnabled("online_booking", "internal-team") {
		t.Error("Internal user harus mendapatkan feature ON")
	}
}

func TestFeatureService_PercentageRollout_10Percent(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 10},
	})

	// Cek user yang hash-nya menghasilkan bucket 0-9 atau 10-99
	// Gunakan tes empiris karena hashing SHA-256 menghasilkan bucket spesifik
	// Cari 1 user di bucket <= 9 dan 1 user di bucket >= 10
	var userOn, userOff string
	for i := 0; i < 1000; i++ {
		userID := fmt.Sprintf("user-%d", i)
		if fs.IsEnabled("online_booking", userID) {
			if userOn == "" {
				userOn = userID
			}
		} else {
			if userOff == "" {
				userOff = userID
			}
		}
		if userOn != "" && userOff != "" {
			break
		}
	}

	if !fs.IsEnabled("online_booking", userOn) {
		t.Errorf("UserID %s harusnya mendapat feature ON", userOn)
	}
	if fs.IsEnabled("online_booking", userOff) {
		t.Errorf("UserID %s harusnya mendapat feature OFF", userOff)
	}
}

func TestFeatureService_DeterministicSameBucket(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 10},
	})

	userID := "consistent-user"
	runs := make([]bool, 10)

	// Setiap request user yang sama harus dapat hasil yang sama
	for i := 0; i < 10; i++ {
		runs[i] = fs.IsEnabled("online_booking", userID)
	}

	first := runs[0]
	for i, got := range runs {
		if got != first {
			t.Errorf("Run %d: inconsistent result. Expected semua %v, got %v", i, first, got)
		}
	}
}

func TestFeatureService_RolloutIncrease_NoRegression(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 10},
	})

	userID := "test-user"

	// Dulu mungkin OFF (tergantung bucket)
	resultBefore := fs.IsEnabled("online_booking", userID)

	// Naik ke 50% (bucket 0-49 = ON)
	fs.SetFlag(Flag{Key: "online_booking", Enabled: true, RolloutPercentage: 50})
	resultAfter := fs.IsEnabled("online_booking", userID)

	// Jika dulu ON, harus tetap ON
	// Jika dulu OFF, bisa jadi ON sekarang (naik ke 50%)
	if resultBefore && !resultAfter {
		t.Error("User yang dulu ON tidak boleh kehilangan feature setelah rollout naik")
	}
}

func TestFeatureService_KillSwitch(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 100},
	})

	// Fitur ON dulu
	if !fs.IsEnabled("online_booking", "user-123") {
		t.Error("Fitur harus ON dengan rollout 100%")
	}

	// Kill switch: Set Enabled = false
	fs.SetFlag(Flag{Key: "online_booking", Enabled: false})

	if fs.IsEnabled("online_booking", "user-123") {
		t.Error("Fitur harus OFF setelah kill switch")
	}
}

func TestFeatureService_UnknownFlag(t *testing.T) {
	fs := NewFeatureService([]Flag{})

	if fs.IsEnabled("unknown_feature", "user-123") {
		t.Error("Unknown flag harus default OFF (safe fallback)")
	}
}

func TestBookingService_FlagOff(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: false},
	})
	metrics := NewMetrics()
	bs := NewBookingService(fs, metrics)

	resp := bs.CreateBooking(BookingRequest{UserID: "user-123"})
	if resp.Flow != "legacy" {
		t.Errorf("Expected flow 'legacy', got '%s'", resp.Flow)
	}
	if !resp.Success {
		t.Error("Expected legacy flow to succeed")
	}
}

func TestBookingService_FlagOn(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 100},
	})
	metrics := NewMetrics()
	bs := NewBookingService(fs, metrics)

	resp := bs.CreateBooking(BookingRequest{UserID: "user-123"})
	if resp.Flow != "online_booking" {
		t.Errorf("Expected flow 'online_booking', got '%s'", resp.Flow)
	}
	if !resp.Success {
		t.Error("Expected online booking flow to succeed")
	}
}

func TestBookingService_KillSwitch_Scenario(t *testing.T) {
	fs := NewFeatureService([]Flag{
		{Key: "online_booking", Enabled: true, RolloutPercentage: 100},
	})
	metrics := NewMetrics()
	bs := NewBookingService(fs, metrics)

	// Skenario: Bug di production! Kita nyalakan simulated failure
	bs.SetSimulateFail(true)

	resp1 := bs.CreateBooking(BookingRequest{UserID: "user-123"})
	if resp1.Flow != "online_booking" || resp1.Success {
		t.Error("Expected online booking to fail when bug is active")
	}

	// KILL SWITCH: feature OFF (tanpa deploy)
	fs.SetFlag(Flag{Key: "online_booking", Enabled: false})

	resp2 := bs.CreateBooking(BookingRequest{UserID: "user-123"})
	if resp2.Flow != "legacy" || !resp2.Success {
		t.Error("Expected legacy fallback to succeed after kill switch")
	}
}

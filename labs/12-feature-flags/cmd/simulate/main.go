package main

import (
	"fmt"
	"math/rand"
	"time"

	featureflags "github.com/lukman-ss/software-engineering-lab/labs/12-feature-flags"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	service := featureflags.NewFeatureService([]featureflags.Flag{
		{Key: "online_booking", Enabled: false},
	})
	metrics := featureflags.NewMetrics()
	bookingService := featureflags.NewBookingService(service, metrics)

	users := generateUsers(100)

	fmt.Println("=== Feature Flag Rollout Simulation ===")
	fmt.Println()

	// Phase 1: OFF
	phase("Phase 1: Feature OFF", service, bookingService, users, metrics)

	// Phase 2: Internal
	service.SetFlag(featureflags.Flag{Key: "online_booking", Enabled: true, InternalUsers: []string{"internal-1", "internal-2", "internal-3"}})
	phase("Phase 2: Internal Rollout", service, bookingService, users, metrics)

	// Phase 3: 10%
	service.SetFlag(featureflags.Flag{Key: "online_booking", Enabled: true, RolloutPercentage: 10, InternalUsers: []string{"internal-1", "internal-2", "internal-3"}})
	phase("Phase 3: 10% Rollout", service, bookingService, users, metrics)

	// Phase 4: 50%
	service.SetFlag(featureflags.Flag{Key: "online_booking", Enabled: true, RolloutPercentage: 50, InternalUsers: []string{"internal-1", "internal-2", "internal-3"}})
	phase("Phase 4: 50% Rollout", service, bookingService, users, metrics)

	// Phase 5: 100%
	service.SetFlag(featureflags.Flag{Key: "online_booking", Enabled: true, RolloutPercentage: 100, InternalUsers: []string{"internal-1", "internal-2", "internal-3"}})
	phase("Phase 5: 100% Rollout", service, bookingService, users, metrics)
	
	// Phase 6: Bug Detected! Turn on kill switch
	fmt.Println("=== INCIDENT IN PRODUCTION ===")
	bookingService.SetSimulateFail(true)
	fmt.Println("Error rate spiked!")
	phase("Phase 6a: 100% Rollout with Bugs", service, bookingService, users, metrics)

	service.SetFlag(featureflags.Flag{Key: "online_booking", Enabled: false})
	phase("Phase 6b: Kill Switch Engaged (Feature OFF)", service, bookingService, users, metrics)
}
func phase(name string, service *featureflags.FeatureService, bookingService *featureflags.BookingService, users []string, metrics *featureflags.Metrics) {
	fmt.Printf("--- %s ---\n", name)

	// Reset metrics
	metrics.Reset()

	// Simulate requests
	for _, user := range users {
		start := time.Now()
		resp := bookingService.CreateBooking(featureflags.BookingRequest{UserID: user, BranchID: "BR-01", MechanicID: "M-01", ServiceType: "OIL"})
		duration := time.Since(start)

		metrics.IncrementTotal()
		metrics.RecordFlow(resp.Flow)
		metrics.RecordBooking(resp.Success)
		metrics.RecordResponseTime(duration)
	}

	snapshot := metrics.Snapshot()
	fmt.Println(snapshot)
	fmt.Println()
}

func generateUsers(n int) []string {
	users := make([]string, n)
	internals := []string{"internal-1", "internal-2", "internal-3"}
	for i := range users {
		if i < 3 {
			users[i] = internals[i]
		} else {
			users[i] = fmt.Sprintf("user-%d", i)
		}
	}
	return users
}
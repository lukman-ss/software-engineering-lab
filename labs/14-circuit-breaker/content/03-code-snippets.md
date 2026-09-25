# Code Snippets

Dokumentasi potongan kode nyata yang disarikan langsung dari implementasi repositori Lab 14 Circuit Breaker.

---

## Snippet 1 — Struktur Definisi dan Konstruktor Circuit Breaker

Source File: `internal/circuitbreaker/circuit_breaker.go`
Purpose: Mendefinisikan struktur data Circuit Breaker, jenis status (State), parameter konfigurasi (Config), dan inisialisasi default yang aman.

```go
type State string

const (
	StateClosed   State = "CLOSED"
	StateOpen     State = "OPEN"
	StateHalfOpen State = "HALF_OPEN"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type Config struct {
	FailureThreshold int
	OpenTimeout      time.Duration
	HalfOpenMaxCalls int
}

type CircuitBreaker struct {
	mu                   sync.Mutex
	state                State
	failureCount         int
	consecutiveSuccesses int
	halfOpenCalls        int
	lastStateChange      time.Time
	config               Config
	now                  func() time.Time
}

func New(cfg Config) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 3
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 5 * time.Second
	}
	if cfg.HalfOpenMaxCalls <= 0 {
		cfg.HalfOpenMaxCalls = 1
	}

	return &CircuitBreaker{
		state:           StateClosed,
		config:          cfg,
		lastStateChange: time.Now(),
		now:             time.Now,
	}
}
```

Explanation:
Kode di atas mendefinisikan tiga status inti (`CLOSED`, `OPEN`, `HALF_OPEN`) dan variabel *error* sentral `ErrCircuitOpen`. Struktur dilindungi oleh `sync.Mutex` untuk menjamin keamanan konkurensi antar *goroutine*. Field `now` diinisialisasi secara fleksibel untuk mempermudah injeksi waktu tiruan (*mock clock*) saat pengujian deterministik.

---

## Snippet 2 — Transisi State Menggunakan Evaluasi Waktu Lazy

Source File: `internal/circuitbreaker/circuit_breaker.go`
Purpose: Mengkalkulasi perpindahan status dari `OPEN` ke `HALF_OPEN` secara instan (*on-demand*) tanpa membebani memori dengan *goroutine timer background*.

```go
func (cb *CircuitBreaker) checkStateTransitionLocked() {
	if cb.state == StateOpen && cb.now().Sub(cb.lastStateChange) >= cb.config.OpenTimeout {
		cb.state = StateHalfOpen
		cb.halfOpenCalls = 0
		cb.consecutiveSuccesses = 0
		cb.lastStateChange = cb.now()
	}
}
```

Explanation:
Fungsi internal ini dipanggil saat kunci mutex telah dipegang (`Locked`). Jika sirkuit sedang berstatus `OPEN` dan waktu jeda istirahat (`OpenTimeout`) telah berlalu, status secara otomatis bertransisi ke `HALF_OPEN` serta mereset penghitung panggilan *probe* (`halfOpenCalls`).

---

## Snippet 3 — Mekanisme Eksekusi dengan Proteksi Mesin Status

Source File: `internal/circuitbreaker/circuit_breaker.go`
Purpose: Mengelola aliran eksekusi fungsi pemanggilan downstream berdasarkan status terkini dari Circuit Breaker.

```go
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	cb.checkStateTransitionLocked()

	switch cb.state {
	case StateOpen:
		cb.mu.Unlock()
		return ErrCircuitOpen

	case StateHalfOpen:
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		cb.halfOpenCalls++
		cb.mu.Unlock()

		err := fn()

		cb.mu.Lock()
		defer cb.mu.Unlock()
		if err != nil {
			cb.state = StateOpen
			cb.lastStateChange = cb.now()
			cb.failureCount = 0
			cb.halfOpenCalls = 0
			return err
		}

		cb.consecutiveSuccesses++
		if cb.consecutiveSuccesses >= cb.config.HalfOpenMaxCalls {
			cb.state = StateClosed
			cb.lastStateChange = cb.now()
			cb.failureCount = 0
			cb.halfOpenCalls = 0
		}
		return nil

	default: // StateClosed
		cb.mu.Unlock()

		err := fn()

		cb.mu.Lock()
		defer cb.mu.Unlock()
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.config.FailureThreshold {
				cb.state = StateOpen
				cb.lastStateChange = cb.now()
			}
			return err
		}

		cb.failureCount = 0
		return nil
	}
}
```

Explanation:
Logika utama proteksi eksekusi:
- Saat `StateOpen`: Menolak panggilan seketika (*fail-fast*) tanpa membuka blok kode fungsi `fn()`.
- Saat `StateHalfOpen`: Mengizinkan pemanggilan pengujian (*probe*) dalam jumlah terbatas (`HalfOpenMaxCalls`). Kegagalan satu kali akan mengembalikan status ke `StateOpen`, sedangkan keberhasilan akan menormalisasi status kembali ke `StateClosed`.
- Saat `StateClosed`: Mengeksekusi fungsi secara normal. Jika terjadi error berulang kali melampaui `FailureThreshold`, status sirkuit langsung diputus menjadi `StateOpen`.

---

## Snippet 4 — Pembungkusan Pemanggilan Klien dengan Circuit Breaker

Source File: `internal/checkout/service.go`
Purpose: Mengintegrasikan Circuit Breaker ke dalam logika domain bisnis layanan checkout saat mengeksekusi pembayaran.

```go
func (s *Service) Checkout(ctx context.Context) error {
	if s.cb != nil {
		err := s.cb.Execute(func() error {
			return s.paymentClient.ProcessPayment(ctx)
		})
		if err != nil {
			return fmt.Errorf("checkout payment failed (with CB): %w", err)
		}
		return nil
	}

	// Without Circuit Breaker
	err := s.paymentClient.ProcessPayment(ctx)
	if err != nil {
		return fmt.Errorf("checkout payment failed (no CB): %w", err)
	}
	return nil
}
```

Explanation:
Metode `Checkout` membungkus pemanggilan jaringan `paymentClient.ProcessPayment(ctx)` di dalam closure `cb.Execute(...)`. Apabila sirkuit terbuka atau downstream bermasalah, error akan diteruskan secara rapi dengan konteks penanganan kegagalan yang tepat.

---

## Snippet 5 — Pengujian Integrasi: Deteksi Gagal dan Pemulihan

Source File: `tests/integration_test.go`
Purpose: Melakukan verifikasi alur pemutusan sirkuit saat dependensi mengalami kegagalan dan transisi pemulihannya ke status semula.

```go
func TestCircuitBreakerIntegration(t *testing.T) {
	fakeServer := payment.NewFakeServer(250 * time.Millisecond)
	defer fakeServer.Close()

	clientTimeout := 50 * time.Millisecond
	paymentClient := payment.NewClient(fakeServer.URL(), clientTimeout)

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 2,
		OpenTimeout:      100 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})

	svc := checkout.NewService(paymentClient, cb)

	t.Run("downstream fails, CB trips open", func(t *testing.T) {
		fakeServer.SetMode(payment.ModeDown)

		// 1st failure
		_ = svc.Checkout(context.Background())
		// 2nd failure trips CB
		_ = svc.Checkout(context.Background())

		if cb.State() != circuitbreaker.StateOpen {
			t.Fatalf("expected OPEN state, got %s", cb.State())
		}

		// 3rd request fails fast
		err := svc.Checkout(context.Background())
		if !errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			t.Fatalf("expected ErrCircuitOpen, got %v", err)
		}
	})

	t.Run("cooldown and recovery", func(t *testing.T) {
		time.Sleep(150 * time.Millisecond)
		fakeServer.SetMode(payment.ModeHealthy)

		if cb.State() != circuitbreaker.StateHalfOpen {
			t.Fatalf("expected HALF_OPEN state, got %s", cb.State())
		}

		err := svc.Checkout(context.Background())
		if err != nil {
			t.Fatalf("probe failed unexpectedly: %v", err)
		}

		if cb.State() != circuitbreaker.StateClosed {
			t.Fatalf("expected CLOSED state after probe, got %s", cb.State())
		}
	})
}
```

Explanation:
Pengujian integrasi ini membuktikan secara empiris:
1. Dua panggilan berturut-turut yang gagal menyebabkan pemutus beralih ke `StateOpen`.
2. Panggilan ketiga seketika menghasilkan error `ErrCircuitOpen` (*fail-fast*).
3. Setelah periode *cooldown* 100ms berlalu dan server dipulihkan ke status `ModeHealthy`, panggilan probe dalam fase `HALF_OPEN` berhasil mengembalikan status ke `StateClosed`.

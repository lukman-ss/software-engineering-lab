# Code Snippets

## Snippet 1 — Model Account dengan Lock Berbasis Channel dan Context Timeout

Source File: internal/bank/account.go
Purpose: Mensimulasikan resource locking dengan kemampuan deteksi timeout (representasi deadlock abort) menggunakan channel Go.

```go
package bank

import (
	"context"
	"errors"
)

var ErrDeadlock = errors.New("deadlock victim")

type Account struct {
	ID      string
	Balance int
	ch      chan struct{}
}

func NewAccount(id string, bal int) *Account {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &Account{ID: id, Balance: bal, ch: ch}
}

func (a *Account) Lock(ctx context.Context) error {
	select {
	case <-a.ch:
		return nil
	case <-ctx.Done():
		return ErrDeadlock
	}
}

func (a *Account) Unlock() {
	a.ch <- struct{}{}
}
```

Explanation:
Model `Account` mengontrol akses konkurensi melalui buffered channel berkapasitas 1. Metode `Lock` mencoba mengambil token dari channel. Apabila token sedang dipegang oleh goroutine lain dan batas waktu context terlampaui (`ctx.Done()`), metode ini mengembalikan error `ErrDeadlock` yang mewakili pemilihan transaksi sebagai *deadlock victim*.

---

## Snippet 2 — Transfer Naif yang Memicu Deadlock (Circular Wait)

Source File: internal/transfer/transfer.go
Purpose: Mengilustrasikan alur penguncian naif berdasarkan urutan parameter pemanggilan yang rentan terhadap circular wait.

```go
// TransferNaive creates a potential deadlock by locking from then to.
func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {
	if err := from.Lock(ctx); err != nil {
		return err
	}
	defer from.Unlock()

	time.Sleep(delay) // simulate transaction duration

	if err := to.Lock(ctx); err != nil {
		return err
	}
	defer to.Unlock()

	from.Balance -= amount
	to.Balance += amount
	return nil
}
```

Explanation:
Operasi ini mengunci akun `from` terlebih dahulu, kemudian mengunci akun `to`. Adanya `time.Sleep(delay)` mensimulasikan waktu pemrosesan di dalam blok transaksi. Ketika dua proses melakukan transfer dengan arah yang berlawanan (A -> B dan B -> A), proses pertama memegang lock A sambil menunggu B, sedangkan proses kedua memegang lock B sambil menunggu A.

---

## Snippet 3 — Pencegahan Deadlock dengan Urutan Penguncian (Lock Ordering)

Source File: internal/transfer/transfer.go
Purpose: Mencegah circular wait dengan selalu mengunci akun dalam urutan alfabetis ID akun.

```go
// TransferOrdered prevents deadlock by always locking accounts in alphabetical order.
func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {
	first, second := acc1, acc2
	if acc1.ID > acc2.ID {
		first, second = acc2, acc1
	}

	if err := first.Lock(ctx); err != nil {
		return err
	}
	defer first.Unlock()

	time.Sleep(delay)

	if err := second.Lock(ctx); err != nil {
		return err
	}
	defer second.Unlock()

	acc1.Balance -= amount
	acc2.Balance += amount
	return nil
}
```

Explanation:
Sebelum mengakuisisi lock, fungsi membandingkan `acc1.ID` dan `acc2.ID` secara leksikografis. Akun dengan ID lebih kecil selalu dikunci sebagai `first`, dan akun lainnya sebagai `second`. Dengan aturan penguncian yang konsisten di semua transaksi, siklus tunggu lingkaran (*circular wait*) tidak dapat terbentuk.

---

## Snippet 4 — Pemulihan Deadlock dengan Pola Retry di Tingkat Aplikasi

Source File: internal/transfer/transfer.go
Purpose: Menangani kesalahan pembatalan akibat deadlock secara transparan dengan perulangan coba ulang.

```go
// TransferWithRetry handles deadlocks by retrying the operation.
func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		// Use a short timeout for the attempt to simulate deadlock monitor aborting fast
		attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		err := TransferNaive(attemptCtx, from, to, amount, delay)
		cancel()

		if err == nil {
			return nil // Success
		}
		if err != bank.ErrDeadlock && err != context.DeadlineExceeded {
			return err // Some other error
		}
		
		// Backoff before retry
		time.Sleep(2 * time.Millisecond)
	}
	return bank.ErrDeadlock
}
```

Explanation:
Fungsi membungkus pemanggilan transaksi dalam perulangan hingga batas `maxRetries`. Jika terjadi `bank.ErrDeadlock` atau `context.DeadlineExceeded`, sistem melakukan penundaan (*backoff*) singkat sebesar 2 milidetik sebelum mencoba kembali transaksi naif dari awal.

---

## Snippet 5 — Verifikasi Otomatis Terjadinya Deadlock

Source File: tests/transfer_test.go
Purpose: Membuktikan secara deterministik bahwa transfer konkuren naif bolak-balik menghasilkan korban deadlock.

```go
func TestDeadlockOccurrence(t *testing.T) {
	accA := bank.NewAccount("A", 100)
	accB := bank.NewAccount("B", 100)

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		errs[0] = transfer.TransferNaive(ctx, accA, accB, 10, 10*time.Millisecond)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		errs[1] = transfer.TransferNaive(ctx, accB, accA, 10, 10*time.Millisecond)
	}()

	wg.Wait()

	hasDeadlock := (errs[0] == bank.ErrDeadlock) || (errs[1] == bank.ErrDeadlock)
	if !hasDeadlock {
		t.Fatalf("expected at least one deadlock victim, got errs: %v, %v", errs[0], errs[1])
	}
}
```

Explanation:
Dua goroutine bersaing mentransfer antara akun A dan B secara simultan dengan durasi transaksi internal 10ms dan batas timeout context 20ms. Pengujian memvalidasi bahwa minimal salah satu goroutine menerima error `bank.ErrDeadlock`.

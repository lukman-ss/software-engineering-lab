package isolation_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
)

// TestSerializable_ConcurrentUpdate_SerializationFailure proves that
// SERIALIZABLE isolation raises 40001 when conflicts are detected.
func TestSerializable_ConcurrentUpdate_SerializationFailure(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

	var wg sync.WaitGroup
	var errA, errB error
	wg.Add(2)

	tx1Read := make(chan struct{})
	tx2Read := make(chan struct{})
	release := make(chan struct{})

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			errA = err
			return
		}
		defer tx.Rollback()

		_, _ = repo.GetBalance(ctx, tx, 1)
		close(tx1Read)
		<-release

		_, errA = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 500000 WHERE id = 1")
		if errA == nil {
			errA = tx.Commit()
		}
	}()

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			errB = err
			return
		}
		defer tx.Rollback()

		_, _ = repo.GetBalance(ctx, tx, 1)
		close(tx2Read)
		<-release

		_, errB = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 500000 WHERE id = 1")
		if errB == nil {
			errB = tx.Commit()
		}
	}()

	<-tx1Read
	<-tx2Read
	close(release)
	wg.Wait()

	t.Logf("Serializable Concurrent Update Error: errA=%v, errB=%v", errA, errB)

	// One should succeed, the other MUST fail with 40001 (serialization failure)
	hasSerializationFailure := (errA != nil && isolation.IsSerializationError(errA)) ||
		(errB != nil && isolation.IsSerializationError(errB))

	if !hasSerializationFailure {
		t.Fatalf("expected serialization failure (40001) in one of the concurrent transactions")
	}
}

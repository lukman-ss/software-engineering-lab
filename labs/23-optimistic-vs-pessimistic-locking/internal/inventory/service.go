package inventory

import (
	"math/rand"
	"time"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (svc *Service) DeductNaive(id int, qty int) error {
	return svc.store.NaiveDeduct(id, qty)
}

func (svc *Service) DeductPessimistic(id int, qty int) error {
	return svc.store.PessimisticDeduct(id, qty)
}

func (svc *Service) DeductOptimisticDirect(id int, qty int) error {
	return svc.store.OptimisticDeduct(id, qty)
}

func (svc *Service) DeductOptimisticWithRetry(id int, qty int, maxRetries int) error {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := svc.store.OptimisticDeduct(id, qty)
		if err == nil {
			return nil
		}
		if err != ErrOptimisticLock {
			return err
		}
		if attempt == maxRetries {
			return ErrOptimisticLock
		}
		// Jittered exponential backoff
		sleepDuration := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond
		time.Sleep(sleepDuration)
	}
	return ErrOptimisticLock
}

func (svc *Service) DeductAtomic(id int, qty int) error {
	return svc.store.AtomicDeduct(id, qty)
}

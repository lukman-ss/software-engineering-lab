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

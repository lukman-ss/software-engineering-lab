package main

import (
	"context"
	"time"
)

type PurchaseOrder struct {
	ID        string
	Amount    float64
	CreatedAt time.Time
}

type Repository interface {
	FetchOrders(ctx context.Context, rangeType string) ([]PurchaseOrder, error)
}

type SlowRepository struct{}

func (r *SlowRepository) FetchOrders(ctx context.Context, rangeType string) ([]PurchaseOrder, error) {
	if rangeType == "1y" {
		time.Sleep(4 * time.Second)
	} else {
		time.Sleep(10 * time.Millisecond)
	}
	return []PurchaseOrder{{ID: "po-1", Amount: 100}}, nil
}

type FastRepository struct{}

func (r *FastRepository) FetchOrders(ctx context.Context, rangeType string) ([]PurchaseOrder, error) {
	time.Sleep(10 * time.Millisecond)
	return []PurchaseOrder{{ID: "po-1", Amount: 100}}, nil
}

type Service struct {
	repo Repository
}

func (s *Service) GetOrders(ctx context.Context, rangeType string) ([]PurchaseOrder, time.Duration, error) {
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	orders, err := s.repo.FetchOrders(ctx, rangeType)
	dur := time.Since(start)
	return orders, dur, err
}

type Handler struct {
	svc *Service
}

type RequestResult struct {
	HandlerDuration    time.Duration
	ServiceDuration    time.Duration
	RepositoryDuration time.Duration
	TotalDuration      time.Duration
}

func (h *Handler) HandleGetOrders(ctx context.Context, rangeType string) (RequestResult, error) {
	totalStart := time.Now()
	handlerStart := time.Now()

	time.Sleep(5 * time.Millisecond)
	handlerPrepDur := time.Since(handlerStart)

	orders, svcDur, err := h.svc.GetOrders(ctx, rangeType)
	if err != nil {
		return RequestResult{}, err
	}
	_ = orders

	totalDur := time.Since(totalStart)
	repoDur := svcDur - 10*time.Millisecond

	return RequestResult{
		HandlerDuration:    handlerPrepDur,
		ServiceDuration:    svcDur,
		RepositoryDuration: repoDur,
		TotalDuration:      totalDur,
	}, nil
}

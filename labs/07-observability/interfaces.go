package observability

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrPDFGeneration   = errors.New("pdf generation failed")
	ErrInventory       = errors.New("inventory reservation failed")
	ErrCommission      = errors.New("commission calculation failed")
	ErrNotification    = errors.New("notification failed")
)

type InvoiceRepository interface {
	Load(ctx context.Context, invoiceID string) error
}

type InventoryService interface {
	Reserve(ctx context.Context, invoiceID string) error
}

type CommissionService interface {
	Calculate(ctx context.Context, invoiceID string) error
}

type PDFGenerator interface {
	Generate(ctx context.Context, invoiceID string) error
}

type NotificationService interface {
	Send(ctx context.Context, invoiceID string) error
}

// ConfigurableDependency simulates I/O operations with configurable duration and errors.
type ConfigurableDependency struct {
	Delay time.Duration
	Err   error
}

func (d ConfigurableDependency) Execute(ctx context.Context) error {
	if d.Delay > 0 {
		select {
		case <-time.After(d.Delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return d.Err
}

type FakeInvoiceRepo struct {
	ConfigurableDependency
}

func (f *FakeInvoiceRepo) Load(ctx context.Context, invoiceID string) error {
	return f.Execute(ctx)
}

type FakeInventoryService struct {
	ConfigurableDependency
}

func (f *FakeInventoryService) Reserve(ctx context.Context, invoiceID string) error {
	return f.Execute(ctx)
}

type FakeCommissionService struct {
	ConfigurableDependency
}

func (f *FakeCommissionService) Calculate(ctx context.Context, invoiceID string) error {
	return f.Execute(ctx)
}

type FakePDFGenerator struct {
	ConfigurableDependency
}

func (f *FakePDFGenerator) Generate(ctx context.Context, invoiceID string) error {
	return f.Execute(ctx)
}

type FakeNotificationService struct {
	ConfigurableDependency
}

func (f *FakeNotificationService) Send(ctx context.Context, invoiceID string) error {
	return f.Execute(ctx)
}

type ScenarioConfig struct {
	RepoDelay         time.Duration
	RepoErr           error
	InventoryDelay    time.Duration
	InventoryErr      error
	CommissionDelay   time.Duration
	CommissionErr     error
	PDFDelay          time.Duration
	PDFErr            error
	NotificationDelay time.Duration
	NotificationErr   error
}

func NewFakeDependencies(cfg ScenarioConfig) (InvoiceRepository, InventoryService, CommissionService, PDFGenerator, NotificationService) {
	return &FakeInvoiceRepo{ConfigurableDependency: ConfigurableDependency{Delay: cfg.RepoDelay, Err: cfg.RepoErr}},
		&FakeInventoryService{ConfigurableDependency: ConfigurableDependency{Delay: cfg.InventoryDelay, Err: cfg.InventoryErr}},
		&FakeCommissionService{ConfigurableDependency: ConfigurableDependency{Delay: cfg.CommissionDelay, Err: cfg.CommissionErr}},
		&FakePDFGenerator{ConfigurableDependency: ConfigurableDependency{Delay: cfg.PDFDelay, Err: cfg.PDFErr}},
		&FakeNotificationService{ConfigurableDependency: ConfigurableDependency{Delay: cfg.NotificationDelay, Err: cfg.NotificationErr}}
}

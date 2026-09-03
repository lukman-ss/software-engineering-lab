package observability

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
)

type UnsafeInvoiceService struct {
	Repo         InvoiceRepository
	Inventory    InventoryService
	Commission   CommissionService
	PDF          PDFGenerator
	Notification NotificationService
	LogWriter    io.Writer
}

func (s *UnsafeInvoiceService) Process(ctx context.Context, invoiceID string) error {
	start := time.Now()

	if err := s.Repo.Load(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := s.Inventory.Reserve(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := s.Commission.Calculate(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := s.PDF.Generate(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := s.Notification.Send(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	duration := time.Since(start)
	s.logf("INFO request completed duration_ms=%d", duration.Milliseconds())
	return nil
}

func (s *UnsafeInvoiceService) logf(format string, args ...any) {
	if s.LogWriter != nil {
		fmt.Fprintf(s.LogWriter, format+"\n", args...)
	}
}

func (s *UnsafeInvoiceService) HTTPHandler(resolveScenario func(r *http.Request) ScenarioConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())
		invoiceID := r.PathValue("id")
		if invoiceID == "" {
			invoiceID = r.URL.Query().Get("id")
		}

		if invoiceID == "" {
			http.Error(w, "missing invoice id", http.StatusBadRequest)
			return
		}

		if resolveScenario != nil {
			cfg := resolveScenario(r)
			repo, inv, comm, pdf, notif := NewFakeDependencies(cfg)
			s.Repo = repo
			s.Inventory = inv
			s.Commission = comm
			s.PDF = pdf
			s.Notification = notif
		}

		if err := s.Process(r.Context(), invoiceID); err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"success","invoice_id":"%s","request_id":"%s"}`, invoiceID, requestID)))
	})
}

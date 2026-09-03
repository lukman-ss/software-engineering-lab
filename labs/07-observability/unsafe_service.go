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
	Deps      Dependencies
	LogWriter io.Writer
}

func (s *UnsafeInvoiceService) Process(ctx context.Context, invoiceID string) error {
	return s.ProcessWithDeps(ctx, invoiceID, s.Deps)
}

func (s *UnsafeInvoiceService) ProcessWithDeps(ctx context.Context, invoiceID string, deps Dependencies) error {
	start := time.Now()

	if err := deps.Repo.Load(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := deps.Inventory.Reserve(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := deps.Commission.Calculate(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := deps.PDF.Generate(ctx, invoiceID); err != nil {
		s.logf("ERROR request failed: %v", err)
		return err
	}

	if err := deps.Notification.Send(ctx, invoiceID); err != nil {
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

		deps := s.Deps
		if resolveScenario != nil {
			cfg := resolveScenario(r)
			deps = NewDependencies(cfg)
		}

		if err := s.ProcessWithDeps(r.Context(), invoiceID, deps); err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"success","invoice_id":"%s","request_id":"%s"}`, invoiceID, requestID)))
	})
}

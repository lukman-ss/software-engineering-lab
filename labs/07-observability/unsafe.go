package observability

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UnsafeInvoiceProcessor processes invoice with plain unstructured/coarse logs.
// It lacks:
// - Request ID / Trace context
// - Dependency breakdown metrics
// - OpenTelemetry trace child spans
// - Granular error attribution
type UnsafeInvoiceProcessor struct {
	Repo         InvoiceRepository
	Inventory    InventoryService
	Commission   CommissionService
	PDF          PDFGenerator
	Notification NotificationService
	LogWriter    io.Writer
}

func (p *UnsafeInvoiceProcessor) ProcessInvoice(ctx context.Context, invoiceID string) error {
	start := time.Now()

	// 1. Database
	if err := p.Repo.Load(ctx, invoiceID); err != nil {
		p.logf("ERROR request failed: %v", err)
		return err
	}

	// 2. Inventory
	if err := p.Inventory.Reserve(ctx, invoiceID); err != nil {
		p.logf("ERROR request failed: %v", err)
		return err
	}

	// 3. Commission
	if err := p.Commission.Calculate(ctx, invoiceID); err != nil {
		p.logf("ERROR request failed: %v", err)
		return err
	}

	// 4. Generate PDF
	if err := p.PDF.Generate(ctx, invoiceID); err != nil {
		p.logf("ERROR request failed: %v", err)
		return err
	}

	// 5. Send Notification
	if err := p.Notification.Send(ctx, invoiceID); err != nil {
		p.logf("ERROR request failed: %v", err)
		return err
	}

	duration := time.Since(start)
	// Coarse log without breakdown
	p.logf("INFO request completed duration_ms=%d", duration.Milliseconds())
	return nil
}

func (p *UnsafeInvoiceProcessor) logf(format string, args ...any) {
	if p.LogWriter != nil {
		fmt.Fprintf(p.LogWriter, format+"\n", args...)
	}
}

// ServeHTTP handles HTTP requests without child spans or granular metrics.
func (p *UnsafeInvoiceProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.URL.Query().Get("id")
	if invoiceID == "" {
		http.Error(w, "missing invoice id", http.StatusBadRequest)
		return
	}

	if err := p.ProcessInvoice(r.Context(), invoiceID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

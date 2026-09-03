package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type contextTraceKey struct{}

// Span represents a single unit of work in a trace.
type Span struct {
	TraceID    string
	SpanID     string
	ParentID   string
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Attributes map[string]string
	Err        error
}

// Tracer provides zero-dependency OpenTelemetry-compatible tracing capabilities.
type Tracer struct {
	mu    sync.RWMutex
	spans []*Span
}

func NewTracer() *Tracer {
	return &Tracer{
		spans: make([]*Span, 0),
	}
}

func (t *Tracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	parent, _ := ctx.Value(contextTraceKey{}).(*Span)

	traceID := ""
	parentID := ""
	if parent != nil {
		traceID = parent.TraceID
		parentID = parent.SpanID
	} else {
		traceID = randomHex(16)
	}

	span := &Span{
		TraceID:    traceID,
		SpanID:     randomHex(8),
		ParentID:   parentID,
		Name:       name,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
	}

	childCtx := context.WithValue(ctx, contextTraceKey{}, span)
	return childCtx, span
}

func (t *Tracer) End(s *Span) {
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)

	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = append(t.spans, s)
}

func (t *Tracer) Spans() []*Span {
	t.mu.RLock()
	defer t.mu.RUnlock()
	res := make([]*Span, len(t.spans))
	copy(res, t.spans)
	return res
}

func (t *Tracer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = make([]*Span, 0)
}

func (s *Span) SetAttribute(key, value string) {
	s.Attributes[key] = value
}

func (s *Span) RecordError(err error) {
	s.Err = err
	if err != nil {
		s.SetAttribute("error", "true")
		s.SetAttribute("error.message", err.Error())
	}
}

func SpanFromContext(ctx context.Context) *Span {
	if s, ok := ctx.Value(contextTraceKey{}).(*Span); ok {
		return s
	}
	return nil
}

func randomHex(bytesLen int) string {
	b := make([]byte, bytesLen)
	rand.Read(b)
	return hex.EncodeToString(b)
}

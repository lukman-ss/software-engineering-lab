package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lukman/software-engineer-lab/pkg/util"
)

type appService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &appService{repo: repo}
}

func (s *appService) NotifyOrderCreated(ctx context.Context, userID, orderID string) error {
	notification := &Notification{
		ID:        util.NewNotificationID(),
		UserID:    userID,
		Type:      TypeOrderCreated,
		Payload:   fmt.Sprintf("{\"order_id\":\"%s\"}", orderID),
		CreatedAt: time.Now(),
	}
	return s.repo.Create(ctx, notification)
}

func (s *appService) NotifyPaymentResult(ctx context.Context, userID, orderID string, success bool) error {
	notifType := TypePaymentSuccess
	if !success {
		notifType = TypePaymentFailed
	}

	notification := &Notification{
		ID:        util.NewNotificationID(),
		UserID:    userID,
		Type:      notifType,
		Payload:   fmt.Sprintf("{\"order_id\":\"%s\",\"success\":%v}", orderID, success),
		CreatedAt: time.Now(),
	}
	return s.repo.Create(ctx, notification)
}
package notifications

import "context"

type NotificationService interface {
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
}

type NotificationImpl struct {
}

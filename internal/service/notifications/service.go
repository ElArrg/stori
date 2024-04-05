package notifications

import "context"

type Service interface {
	SendNotification(ctx context.Context, accountID string, operationName string, payload map[string]any) []error
}

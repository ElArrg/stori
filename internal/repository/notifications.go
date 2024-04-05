package repository

import (
	"context"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Notifications interface {
	GetActiveTemplatesByOperationAndChannels(ctx context.Context, operation string, channels []models.Channel) ([]models.Template, error)
	GetEnabledChannelsByAccountID(ctx context.Context, accountID string) ([]models.Channel, error)
}

package postgres

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/elarrg/stori/ledger/internal/models"
)

type NotificationsRepository struct {
	db *bun.DB
}

func NewNotificationsRepository(db *bun.DB) *NotificationsRepository {
	return &NotificationsRepository{
		db: db,
	}
}

func (n *NotificationsRepository) GetEnabledChannelsByAccountID(ctx context.Context, accountID string) ([]models.Channel, error) {
	activeChannels := make([]models.Channel, 0)

	err := n.db.NewSelect().
		Model((*models.NotificationsSettings)(nil)).
		Column("channel").
		Where("account_id = ?", accountID).
		Where("enabled = true").
		Scan(ctx, &activeChannels)

	if err != nil {
		return nil, err
	}

	return activeChannels, nil
}

func (n *NotificationsRepository) GetActiveTemplatesByOperationAndChannels(ctx context.Context, operation string, channels []models.Channel) ([]models.Template, error) {
	template := make([]models.Template, 0)

	err := n.db.NewSelect().
		Model(&template).
		Where("active = true").
		Where("operation = ?", operation).
		Where("channel IN (?)", bun.In(channels)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return template, nil
}

package repository

import (
	"context"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Accounts interface {
	GetByID(ctx context.Context, id string) (*models.Account, error)
}

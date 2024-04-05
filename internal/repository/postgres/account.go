package postgres

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/elarrg/stori/ledger/internal/models"
)

type AccountRepository struct {
	db *bun.DB
}

func NewAccountRepository(db *bun.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	account := new(models.Account)

	err := a.db.NewSelect().
		Model(account).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return account, nil
}

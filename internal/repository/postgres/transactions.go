package postgres

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/elarrg/stori/ledger/internal/models"
	"github.com/elarrg/stori/ledger/internal/repository"
)

type TransactionRepository struct {
	db *bun.DB
}

func NewTransactionRepository(db *bun.DB) repository.Transactions {
	return &TransactionRepository{
		db: db,
	}
}

func (t *TransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountId string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := t.db.NewSelect().
		Model(&transactions).
		Where("account_id = ?", accountId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionRepository) GetBalanceReportByAccountIDAndType(ctx context.Context, accountID string, balanceType string) (*models.BalanceReport, error) {
	balanceReport := new(models.BalanceReport)

	err := t.db.NewSelect().
		Model(balanceReport).
		ModelTableExpr("transactions as t").
		ColumnExpr("SUM(t.amount) as total_balance").
		ColumnExpr("AVG(t.amount) as average_amount").
		ColumnExpr("? as balance_type", balanceType).
		Where("t.type = ?", balanceType).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return balanceReport, nil
}

func (t *TransactionRepository) GetTransactionsByAccountIDGroupedByMonth(ctx context.Context, accountID string) ([]models.MonthCount, error) {
	monthCount := make([]models.MonthCount, 0)

	err := t.db.NewSelect().
		Model((*models.Transaction)(nil)).
		Column("year", "month").
		ColumnExpr("COUNT(*)").
		Group("year", "month").
		Order("year DESC", "month DESC").
		Scan(ctx, &monthCount)

	if err != nil {
		return nil, err
	}

	return monthCount, nil
}

func (t *TransactionRepository) InsertTransactionsInBulk(ctx context.Context, transaction []models.Transaction) error {
	_, err := t.db.NewInsert().
		Model(&transaction).
		Exec(ctx)

	return err
}

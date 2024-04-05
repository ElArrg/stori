package repository

import (
	"context"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Transactions interface {
	GetTransactionsByAccountID(ctx context.Context, accountId string) ([]models.Transaction, error)
	GetBalanceReportByAccountIDAndType(ctx context.Context, accountID string, balanceType string) (*models.BalanceReport, error)
	GetTransactionsByAccountIDGroupedByMonth(ctx context.Context, accountID string) ([]models.MonthCount, error)

	InsertTransactionsInBulk(ctx context.Context, transaction []models.Transaction) error
}

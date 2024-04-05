package transactions

import (
	"context"
	"io"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Service interface {
	ProcessTransactionsFile(ctx context.Context, reader io.Reader) (summaries []models.BalanceSummary, errs []error)
}

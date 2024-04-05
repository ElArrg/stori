package dispatchers

import (
	"context"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, account models.Account, template models.Template, payload map[string]any) error
}

type Operation string

const (
	AccountSummaryOp = "account-summary"
)

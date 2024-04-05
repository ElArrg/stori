package parser

import (
	"context"
	"io"

	"github.com/elarrg/stori/ledger/internal/models"
)

type Parser interface {
	Parse(ctx context.Context, r io.Reader) ([]models.Transaction, error)
}

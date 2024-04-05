package transactions

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"

	"github.com/elarrg/stori/ledger/internal/models"
	"github.com/elarrg/stori/ledger/internal/repository"
	"github.com/elarrg/stori/ledger/internal/service/notifications"
	"github.com/elarrg/stori/ledger/internal/service/notifications/dispatchers"
	"github.com/elarrg/stori/ledger/internal/service/transactions/parser"
)

type DefaultService struct {
	transRepo  repository.Transactions
	fileParser parser.Parser
	notifSvc   notifications.Service
}

func NewDefaultService(tr repository.Transactions, fp parser.Parser, ns notifications.Service) *DefaultService {
	return &DefaultService{
		transRepo:  tr,
		fileParser: fp,
		notifSvc:   ns,
	}
}

func (d *DefaultService) ProcessTransactionsFile(ctx context.Context, reader io.Reader) (summaries []models.BalanceSummary, errs []error) {
	txns, err := d.fileParser.Parse(ctx, reader)
	if err != nil {
		// todo log
		return nil, append(errs, err)
	}

	err = d.transRepo.InsertTransactionsInBulk(ctx, txns)
	if err != nil {
		// todo log
		return nil, append(errs, errors.New("couldn't store the transactions from the file"))
	}

	// Get unique accounts from transactions
	accountsSet := make(map[string]bool)
	for _, txn := range txns {
		if !accountsSet[txn.AccountID] {
			accountsSet[txn.AccountID] = true
		}
	}

	for accountID, _ := range accountsSet {
		var creditReport, debitReport *models.BalanceReport
		var monthsCount []models.MonthCount

		creditReport, err = d.transRepo.GetBalanceReportByAccountIDAndType(ctx, accountID, models.CreditTransactionType)
		if err != nil {
			// todo log
			errs = append(errs, fmt.Errorf("couldn't get the credit report for account %v", accountID))
			continue
		}

		debitReport, err = d.transRepo.GetBalanceReportByAccountIDAndType(ctx, accountID, models.DebitTransactionType)
		if err != nil {
			// todo log
			errs = append(errs, fmt.Errorf("couldn't get the debit report for account %v", accountID))
			continue
		}

		monthsCount, err = d.transRepo.GetTransactionsByAccountIDGroupedByMonth(ctx, accountID)
		if err != nil {
			// todo log
			errs = append(errs, fmt.Errorf("couldn't get the transactions per month for account %v", accountID))
			continue
		}

		totalBalance := creditReport.TotalBalance + debitReport.TotalBalance

		summary := models.BalanceSummary{
			AccountID:           accountID,
			TotalBalance:        float64(totalBalance) / 100,
			AverageCredit:       creditReport.AverageAmount,
			AverageDebit:        debitReport.AverageAmount,
			TransactionsByMonth: monthsCount,
		}

		summaries = append(summaries, summary)

		// TODO: Publish Events
		payload := make(map[string]any)
		err = mapstructure.Decode(summary, &payload)
		if err != nil {
			errs = append(errs, fmt.Errorf("couldn't encode balance summary for account %v", accountID))
			continue
		}

		e := d.notifSvc.SendNotification(ctx, accountID, dispatchers.AccountSummaryOp, payload)
		if e != nil {
			errs = append(errs, e...)
		}

	}

	return summaries, errs
}

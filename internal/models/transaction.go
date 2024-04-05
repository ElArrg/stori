package models

import "time"

const (
	// DebitTransactionType is the type for debit transactions.
	DebitTransactionType string = "debit"

	// CreditTransactionType is the type for credit transactions.
	CreditTransactionType string = "credit"
)

// Transaction is a financial transaction with its details.
type Transaction struct {
	ID        string     // ID the identifier of this transaction
	AccountID string     // AccountID is the account affected by this transaction.
	Date      time.Time  // Date is when the transaction occurred.
	Amount    int64      // Amount of the transaction, it's represented in cents
	Type      string     // Type of transaction, determines how will affect the balance if as credit or debit operation.
	Year      int        // Year when the transaction occurred
	Month     time.Month // Month when the transaction occurred
}

// BalanceReport represents a report of the balance for a specific type
// of transactions.
type BalanceReport struct {
	AccountID     string  // AccountID related to the report
	TotalBalance  int64   // TotalBalance is the sum of the transactions of type BalanceType
	AverageAmount float64 // TransactionsCount is the total transactions for this report
	BalanceType   string  // BalanceType is the type of this balance report. can be one of DebitTransactionType, CreditTransactionType
}

type MonthCount struct {
	Month time.Month `json:"month"`
	Year  int        `json:"year"`
	Count int64      `json:"count"`
}

type BalanceSummary struct {
	AccountID           string       `mapstructure:"-"`
	TotalBalance        float64      `mapstructure:"totalBalance"`
	AverageCredit       float64      `mapstructure:"averageCredit"`
	AverageDebit        float64      `mapstructure:"averageDebit"`
	TransactionsByMonth []MonthCount `mapstructure:"transactionsByMonth"`
}

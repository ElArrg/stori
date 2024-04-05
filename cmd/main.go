package main

import (
	"context"
	"log"
	"time"

	"github.com/elarrg/stori/ledger/internal/adapters/clients/sendgrid"
	"github.com/elarrg/stori/ledger/internal/repository/postgres"
	"github.com/elarrg/stori/ledger/internal/service/notifications"
	"github.com/elarrg/stori/ledger/internal/service/notifications/dispatchers"

	"github.com/elarrg/stori/ledger/internal/adapters/db"
	"github.com/elarrg/stori/ledger/internal/service/sources"
	"github.com/elarrg/stori/ledger/internal/service/transactions"
	"github.com/elarrg/stori/ledger/internal/service/transactions/parser"
)

func main() {
	config := struct {
		SengridClient sendgrid.ClientConfigs
		PostgresDB    db.PostgresConfig
	}{
		SengridClient: sendgrid.ClientConfigs{
			SenderEmail: "notifications@elarrg.com",
			SenderName:  "ElArrg Notifications",
			Key:         "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",
			Host:        "https://api.sendgrid.com",
			SandboxMode: false,
		},
		PostgresDB: db.PostgresConfig{
			DSN:        "postgresql://postgres:pass-ledger@localhost:5432/ledger",
			QueryDebug: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	filepath := "/Users/elias/Development/interview/stori/ledger/resources/transactions/1_txns.csv"

	// dependency injection
	postgresDB, err := db.NewPostgresDB(&config.PostgresDB)
	if err != nil {
		log.Fatalf("couldn't connect to DB: %v", err)
	}

	// Repositories
	accountRepo := postgres.NewAccountRepository(postgresDB.DB)
	transRepo := postgres.NewTransactionRepository(postgresDB.DB)
	notifRepo := postgres.NewNotificationsRepository(postgresDB.DB)

	// clients
	sendgridClient := sendgrid.NewDefaultClient(&config.SengridClient)

	diskSrcOp := sources.NewDiskOpener()
	csvParser := parser.NewCSVParser([]string{})

	// Services
	emailDispatcher := dispatchers.NewEmailProcessor(sendgridClient)

	notifSvc := notifications.NewDefaultService(notifRepo, accountRepo,
		notifications.WithEmailDispatcher(emailDispatcher),
	)
	transSvc := transactions.NewDefaultService(transRepo, csvParser, notifSvc)

	// Start the process
	file, err := diskSrcOp.OpenFromSource(filepath)
	if err != nil {
		log.Fatalf("couldn't open file from source")
	}

	_, errs := transSvc.ProcessTransactionsFile(ctx, file)
	if len(errs) != 0 {
		log.Printf("Errors while processing the file: %v\n", errs)
	}

}

package main

import (
	"context"
	"log"
	"time"

	"github.com/elarrg/stori/ledger/configs"
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
	conf, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// dependency injection
	postgresDB, err := db.NewPostgresDB(&conf.PostgresDB)
	if err != nil {
		log.Fatalf("couldn't connect to DB: %v", err)
	}

	// Repositories
	accountRepo := postgres.NewAccountRepository(postgresDB.DB)
	transRepo := postgres.NewTransactionRepository(postgresDB.DB)
	notifRepo := postgres.NewNotificationsRepository(postgresDB.DB)

	// clients
	sendgridClient := sendgrid.NewDefaultClient(&conf.Sendgrid)

	diskSrcOp := sources.NewDiskOpener()
	csvParser := parser.NewCSVParser([]string{})

	// Services
	emailDispatcher := dispatchers.NewEmailProcessor(sendgridClient)

	notifSvc := notifications.NewDefaultService(notifRepo, accountRepo,
		notifications.WithEmailDispatcher(emailDispatcher),
	)
	transSvc := transactions.NewDefaultService(transRepo, csvParser, notifSvc)

	// Start the process
	file, err := diskSrcOp.OpenFromSource(conf.Transactions.SourcePath)
	if err != nil {
		log.Fatalf("couldn't open file from source")
	}

	_, errs := transSvc.ProcessTransactionsFile(ctx, file)
	if len(errs) != 0 {
		log.Printf("Errors while processing the file: %v\n", errs)
	}

}

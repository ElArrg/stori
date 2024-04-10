package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresConfig struct {
	DSN        string `koanf:"dsn"`
	QueryDebug bool   `koanf:"query_debug"`
}

type PostgresDB struct {
	DB      *bun.DB
	configs *PostgresConfig
}

func NewPostgresDB(configs *PostgresConfig) (*PostgresDB, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithDSN(configs.DSN),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(pgconn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := sqldb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if configs.QueryDebug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		),
		)
	}

	return &PostgresDB{
		DB: db,
	}, nil
}

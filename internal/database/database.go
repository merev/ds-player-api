package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	const enablePgcrypto = `CREATE EXTENSION IF NOT EXISTS pgcrypto;`

	const playersTable = `
CREATE TABLE IF NOT EXISTS players (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

	if _, err := db.Exec(ctx, enablePgcrypto); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, playersTable); err != nil {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

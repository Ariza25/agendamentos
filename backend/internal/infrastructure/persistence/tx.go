package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type PostgresTxRunner struct {
	db *sql.DB
}

func NewPostgresTxRunner(db *sql.DB) *PostgresTxRunner {
	return &PostgresTxRunner{db: db}
}

func (r *PostgresTxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if r.db == nil {
		return fn(ctx)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

type NopTxRunner struct{}

func (NopTxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

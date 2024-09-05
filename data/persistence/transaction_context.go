package persistence

import (
	"context"
	"database/sql"
	"errors"
)

type transactionContextType string

const transactionContextKey transactionContextType = "persistence.tx"

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TransactionSQL struct {
	Tx *sql.Tx
}

var _ Transaction = (*TransactionSQL)(nil)

func (t TransactionSQL) Commit(_ context.Context) error {
	return t.Tx.Commit()
}

func (t TransactionSQL) Rollback(_ context.Context) error {
	return t.Tx.Rollback()
}

func GetTxFromContext(ctx context.Context) (Transaction, error) {
	tx, ok := ctx.Value(transactionContextKey).(Transaction)
	if !ok {
		return nil, ErrTxContextNotFound
	}
	return tx, nil
}

func CloseTransaction(ctx context.Context, srcErr error) error {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return errors.Join(srcErr, err)
	}
	if recovered := recover(); recovered != nil {
		return errors.Join(srcErr, tx.Rollback(ctx))
	}
	if srcErr != nil {
		errRollback := tx.Rollback(ctx)
		return errors.Join(srcErr, errRollback)
	}
	return tx.Commit(ctx)
}

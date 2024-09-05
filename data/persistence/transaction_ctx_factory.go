package persistence

import (
	"context"
	"database/sql"
)

type TransactionContextFactory interface {
	NewContext(parent context.Context) (context.Context, error)
}

type TransactionContextFactorySQL struct {
	DB     ClientSQL
	Config ConfigTransactionManagerSQL
}

var _ TransactionContextFactory = (*TransactionContextFactorySQL)(nil)

func NewTransactionContextFactorySQL(db ClientSQL, cfg ConfigTransactionManagerSQL) TransactionContextFactorySQL {
	return TransactionContextFactorySQL{
		DB:     db,
		Config: cfg,
	}
}

func (t TransactionContextFactorySQL) NewContext(parent context.Context) (context.Context, error) {
	_, err := GetTxFromContext(parent)
	if err == nil {
		return parent, nil // re-use ctx
	}

	tx, err := t.DB.BeginTx(parent, &sql.TxOptions{
		Isolation: sql.IsolationLevel(t.Config.IsolationLevel),
		ReadOnly:  t.Config.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	return context.WithValue(parent, transactionContextKey, TransactionSQL{Tx: tx}), nil
}

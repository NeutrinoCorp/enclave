package persistence

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/neutrinocorp/geck/logging"
)

type ClientSQL interface {
	PingContext(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Driver() driver.Driver
}

type TransactionalClientSQL struct {
	TransactionContextFactory TransactionContextFactorySQL
	Logger                    logging.Logger
	Next                      ClientSQL
}

var _ ClientSQL = (*TransactionalClientSQL)(nil)

func NewTransactionalClientSQL(factory TransactionContextFactorySQL, logger logging.Logger, next ClientSQL) TransactionalClientSQL {
	return TransactionalClientSQL{
		TransactionContextFactory: factory,
		Logger:                    logger,
		Next:                      next,
	}
}

func (t TransactionalClientSQL) PingContext(ctx context.Context) error {
	return t.Next.PingContext(ctx)
}

func (t TransactionalClientSQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	txRaw, err := GetTxFromContext(ctx)
	if err != nil {
		t.Logger.Warn().WithField("err", err).Write("error getting transaction, using fallback client")
		return t.Next.ExecContext(ctx, query, args...)
	}

	tx, ok := txRaw.(TransactionSQL)
	if !ok {
		t.Logger.Error().Write("error casting transaction structure, using fallback client")
		return t.Next.ExecContext(ctx, query, args...)
	}
	return tx.Tx.ExecContext(ctx, query, args...)
}

func (t TransactionalClientSQL) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	txRaw, err := GetTxFromContext(ctx)
	if err != nil {
		t.Logger.Warn().WithField("err", err).Write("error getting transaction, using fallback client")
		return t.Next.PrepareContext(ctx, query)
	}

	tx, ok := txRaw.(TransactionSQL)
	if !ok {
		t.Logger.Error().Write("error casting transaction structure, using fallback client")
		return t.Next.PrepareContext(ctx, query)
	}
	return tx.Tx.PrepareContext(ctx, query)
}

func (t TransactionalClientSQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	txRaw, err := GetTxFromContext(ctx)
	if err != nil {
		t.Logger.Warn().WithField("err", err).Write("error getting transaction, using fallback client")
		return t.Next.QueryContext(ctx, query, args...)
	}

	tx, ok := txRaw.(TransactionSQL)
	if !ok {
		t.Logger.Error().Write("error casting transaction structure, using fallback client")
		return t.Next.QueryContext(ctx, query, args...)
	}
	return tx.Tx.QueryContext(ctx, query, args...)
}

func (t TransactionalClientSQL) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	txRaw, err := GetTxFromContext(ctx)
	if err != nil {
		t.Logger.Warn().WithField("err", err).Write("error getting transaction, using fallback client")
		return t.Next.QueryRowContext(ctx, query, args...)
	}

	tx, ok := txRaw.(TransactionSQL)
	if !ok {
		t.Logger.Error().Write("error casting transaction structure, using fallback client")
		return t.Next.QueryRowContext(ctx, query, args...)
	}
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

func (t TransactionalClientSQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return t.Next.BeginTx(ctx, opts)
}

func (t TransactionalClientSQL) Driver() driver.Driver {
	return t.Next.Driver()
}

type StatementLoggerClientSQL struct {
	Logger logging.Logger
	Next   ClientSQL
}

var _ ClientSQL = (*StatementLoggerClientSQL)(nil)

func NewStatementLoggerClientSQL(logger logging.Logger, db ClientSQL) StatementLoggerClientSQL {
	return StatementLoggerClientSQL{
		Logger: logger,
		Next:   db,
	}
}

func (l StatementLoggerClientSQL) PingContext(ctx context.Context) error {
	l.Logger.Debug().Write("pinging database")
	return l.Next.PingContext(ctx)
}

func (l StatementLoggerClientSQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	l.Logger.Debug().WithField("statement", query).
		WithField("total_args", len(args)).
		Write("executing query")
	return l.Next.ExecContext(ctx, query, args...)
}

func (l StatementLoggerClientSQL) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	l.Logger.Debug().WithField("statement", query).Write("preparing query")
	return l.Next.PrepareContext(ctx, query)
}

func (l StatementLoggerClientSQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	l.Logger.Debug().WithField("statement", query).
		WithField("total_args", len(args)).
		Write("querying statement")
	return l.Next.QueryContext(ctx, query, args...)
}

func (l StatementLoggerClientSQL) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	l.Logger.Debug().WithField("statement", query).
		WithField("total_args", len(args)).
		Write("querying statement")
	return l.Next.QueryRowContext(ctx, query, args...)
}

func (l StatementLoggerClientSQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	isolationLvl := 0
	readOnly := false
	if opts != nil {
		isolationLvl = int(opts.Isolation)
		readOnly = opts.ReadOnly
	}
	l.Logger.Debug().
		WithField("isolation_level", isolationLvl).
		WithField("read_only", readOnly).
		Write("starting transaction")
	return l.Next.BeginTx(ctx, opts)
}

func (l StatementLoggerClientSQL) Driver() driver.Driver {
	return l.Next.Driver()
}

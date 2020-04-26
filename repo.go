package zergrepo

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

// Repo The wrapper around *sql.DB.
// Provides a number of convenient methods for starting a transaction
// and starting functions and wrapping returnable errors.
type Repo struct {
	db  *sql.DB
	log *zap.Logger
}

// New return new instance Repo.
func New(db *sql.DB, log *zap.Logger) *Repo {
	return &Repo{
		db:  db,
		log: log,
	}
}

// Tx automatically starts a transaction according to the parameters.
// If the options are not sent, the transaction will start with default parameters.
// If the callback returns the error, it will be wrapped and enriched with
// information about where the transaction was called from.
// TODO: Add collect metric.
func (r *Repo) Tx(ctx context.Context, fn func(*sql.Tx) error, opts ...TxOption) error {
	methodName := callerMethodName()

	txOption := &sql.TxOptions{}
	for i := range opts {
		opts[i](txOption)
	}

	tx, err := r.db.BeginTx(ctx, txOption)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	err = fn(tx)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = fmt.Errorf("roolback err: %w", errRollback)
		}

		return fmt.Errorf("%s: %w", methodName, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	return nil
}

// Do a wrapper for database requests.
// If the callback returns the error, it will be wrapped and enriched with
// information about where the transaction was called from.
// TODO: Add collect metric.
func (r *Repo) Do(fn func(*sql.DB) error) error {
	methodName := callerMethodName()

	err := fn(r.db)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	return nil
}

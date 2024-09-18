package pg

import (
	"github.com/jackc/pgx/v5"
)

const (
	stmtTimeoutOp   = "SET local statement_timeout = $1"
	stmtTimeoutTx   = "SET local idle_in_transaction_session_timeout = $1"
	stmtTimeoutLock = "SET local lock_timeout = $1"
)

type PgTransactable interface {
	Transact(PgOp) error
}

type PgTransaction struct {
	PgTransactable
	hnd  PgCommands
	opts PgOptions
}

func (tx *PgTransaction) Transact(fn PgOp) error {
	err := pgx.BeginFunc(tx.opts.Ctx, tx.hnd, func(ptx pgx.Tx) error {
		if err := applyTimeout(ptx, tx.opts); err != nil {
			return err
		}

		return fn(&PgOperation{
			hnd:  ptx,
			opts: tx.opts,
		})
	})

	return err
}

func applyTimeout(tx pgx.Tx, opts PgOptions) error {
	ctx := opts.Ctx

	if opts.TxTimeout > 0 {
		if _, err := tx.Exec(ctx, stmtTimeoutTx, opts.TxTimeout); err != nil {
			return err
		}
	}

	if opts.LockTimeout > 0 {
		if _, err := tx.Exec(ctx, stmtTimeoutLock, opts.LockTimeout); err != nil {
			return err
		}
	}

	if opts.StmtTimeout > 0 {
		if _, err := tx.Exec(ctx, stmtTimeoutOp, opts.StmtTimeout); err != nil {
			return err
		}
	}

	return nil
}

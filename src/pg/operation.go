package pg

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgOperable interface {
	Exec(string, ...interface{}) (uint64, error)
	ExecRaw(string, ...any) (pgconn.CommandTag, error)

	Query(interface{}, string, ...interface{}) error
	QueryRaw(string, ...interface{}) (pgx.Rows, error)

	Get(interface{}, string, ...interface{}) error
}

type PgOperation struct {
	PgOperable
	hnd  PgCommands
	opts PgOptions
}

func (op *PgOperation) Exec(sql string, args ...interface{}) (uint64, error) {
	cmd, err := op.hnd.Exec(op.opts.Ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return uint64(cmd.RowsAffected()), nil
}

func (op *PgOperation) ExecRaw(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return op.hnd.Exec(op.opts.Ctx, sql, args...)
}

func (op *PgOperation) Query(target interface{}, sql string, args ...interface{}) error {
	rows, err := op.hnd.Query(op.opts.Ctx, sql, args...)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(target, rows); err != nil {
		return err
	}

	return nil
}

func (op *PgOperation) QueryRaw(sql string, args ...interface{}) (pgx.Rows, error) {
	return op.hnd.Query(op.opts.Ctx, sql, args...)
}

func (op *PgOperation) Get(target interface{}, sql string, args ...interface{}) error {
	rows, err := op.hnd.Query(op.opts.Ctx, sql, args...)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(target, rows); err != nil {
		return err
	}

	return nil
}

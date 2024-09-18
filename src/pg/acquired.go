package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgAcquirable func(*pgxpool.Conn) error

type PgAcquired struct {
	pool context.Context
	opts PgOptions
}

var (
	pgAcqCtx = struct{ name string }{name: "pgAcquired"}
)

func NewAcquire(pool *pgxpool.Pool, options PgOptions) *PgAcquired {
	ctx := context.WithValue(options.Ctx, pgAcqCtx, pool)
	return &PgAcquired{
		pool: ctx,
		opts: options,
	}
}

func (p *PgAcquired) One() (*pgxpool.Conn, error) {
	pool, ok := p.pool.Value(pgAcqCtx).(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("")
	}

	return pool.Acquire(p.opts.Ctx)
}

func (p *PgAcquired) With(fn PgAcquirable) error {
	pool, ok := p.pool.Value(pgAcqCtx).(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("")
	}

	return pool.AcquireFunc(p.opts.Ctx, fn)
}

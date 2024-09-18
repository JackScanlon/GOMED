package pg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgInitTimeout = 30 * time.Second
)

type Driver struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

var (
	instance *Driver = nil
	lock             = &sync.Mutex{}
)

func GetDB(ctx context.Context) (*Driver, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		ctx, cancel := context.WithTimeout(ctx, pgInitTimeout)
		defer cancel()

		pool, err := pgxpool.New(ctx, buildConnectionString())
		if err != nil {
			return nil, err
		}

		if err := pool.Ping(ctx); err != nil {
			return nil, err
		}

		instance = &Driver{
			pool: pool,
			ctx:  context.Background(),
		}

		return instance, nil
	}

	return instance, nil
}

func (d *Driver) GetPool() *pgxpool.Pool {
	return d.pool
}

func (d *Driver) Ping(options ...PgOption) error {
	opts := d.getOptions(options...)
	return d.pool.Ping(opts.Ctx)
}

func (d *Driver) Acquire(options ...PgOption) *PgAcquired {
	return NewAcquire(d.pool, d.getOptions(options...))
}

func (d *Driver) Stmt(options ...PgOption) *PgOperation {
	return &PgOperation{
		hnd:  d.pool,
		opts: d.getOptions(options...),
	}
}

func (d *Driver) Tx(options ...PgOption) *PgTransaction {
	return &PgTransaction{
		hnd:  d.pool,
		opts: d.getOptions(options...),
	}
}

func (d *Driver) Close() {
	lock.Lock()
	defer lock.Unlock()

	if instance != nil {
		d.GetPool().Close()
		instance = nil
	}
}

func (d *Driver) getOptions(options ...PgOption) PgOptions {
	var opts PgOptions
	for _, opt := range options {
		opt(&opts)
	}

	if opts.Ctx == nil {
		opts.Ctx = d.ctx
	}

	return opts
}

func buildConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		Config.PostgresUsername,
		Config.PostgresPassword,
		Config.PostgresHost,
		Config.PostgresPort,
		Config.PostgresDatabase,
	)
}

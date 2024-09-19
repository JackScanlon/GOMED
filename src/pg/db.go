package pg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
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

func (d *Driver) CreateTableFrom(schema string, name string, obj interface{}, options ...PgOption) error {
	exists, err := d.Exists(schema, name, options...)
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("invalid arguments: table of name '%s' in schema '%s' already exists", schema, name)
	}

	content, err := BuildCreateString(schema, name, obj)
	if err != nil {
		return err
	}

	_, err = d.
		Stmt(options...).
		Exec(content)

	return err
}

func (d *Driver) Exists(schema string, name string, options ...PgOption) (bool, error) {
	var exists bool
	err := d.
		Stmt(options...).
		Get(
			&exists,
			"SELECT EXISTS(SELECT table_name FROM information_schema.tables WHERE table_schema=$1 and table_name=$2)",
			schema, name,
		)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (d *Driver) DropIfExists(schema string, name string, options ...PgOption) error {
	_, err := d.
		Stmt(options...).
		Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", pgx.Identifier{schema, name}.Sanitize()))

	return err
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

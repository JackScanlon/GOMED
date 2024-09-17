package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbTimeout = 30 * time.Second
)

var (
	pool *pgxpool.Pool = nil
)

func GetDb(ctx context.Context) (*pgxpool.Pool, error) {
	if pool == nil {
		var err error
		ctx, cancel := context.WithTimeout(ctx, dbTimeout)
		defer cancel()

		pool, err = pgxpool.New(ctx, buildConnectionString())
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
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

package pg

import (
	"context"
	"time"
)

type PgOptions struct {
	Ctx         context.Context
	TxTimeout   uint64
	StmtTimeout uint64
	LockTimeout uint64
}

type PgOption func(*PgOptions)

func WithStmtContext(ctx context.Context) PgOption {
	return func(pgo *PgOptions) {
		pgo.Ctx = ctx
	}
}

func WithStmtTimeout(t time.Duration) PgOption {
	return func(pgo *PgOptions) {
		pgo.StmtTimeout = uint64(t.Milliseconds())
	}
}

func WithTxTimeout(t time.Duration) PgOption {
	return func(pgo *PgOptions) {
		pgo.TxTimeout = uint64(t.Milliseconds())
	}
}

func WithLockTimeout(t time.Duration) PgOption {
	return func(pgo *PgOptions) {
		pgo.LockTimeout = uint64(t.Milliseconds())
	}
}

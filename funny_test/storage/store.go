package storage

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	dbpool  *pgxpool.Pool
	log     *logrus.Entry
	CloseFn func()
}

func New(ctx context.Context, connConf *pgxpool.Config, log *logrus.Entry) (*Store, error) {
	dbpool, err := pgxpool.ConnectConfig(ctx, connConf)

	if err != nil {
		return nil, err
	}

	return &Store{
		dbpool: dbpool,
		log:    log,

		CloseFn: dbpool.Close,
	}, nil
}

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close()
}

func (s *Store) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return s.dbpool.Query(ctx, sql, args...)
}

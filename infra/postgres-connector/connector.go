package pgconnector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(pgconn string, attempts int, attemptDelay time.Duration, ctxTimeout time.Duration) (*pgxpool.Pool, error) {
	if pgconn == "" {
		return nil, errors.New("empty postgres connection string")
	}
	connConfig, err := pgxpool.ParseConfig(pgconn)
	if err != nil {
		return nil, fmt.Errorf("err parsing postgres connection string: %w", err)
	}

	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		pool, err := pgxpool.ConnectConfig(ctx, connConfig)
		if err == nil {
			return pool, nil
		}
		time.Sleep(attemptDelay)
		attempts--
	}
	return nil, err
}

func ConnectConfig(config *pgxpool.Config, attempts int, attemptDelay time.Duration, ctxTimeout time.Duration) (*pgxpool.Pool, error) {
	var err error
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		pool, err := pgxpool.ConnectConfig(ctx, config)
		if err == nil {
			return pool, nil
		}
		time.Sleep(attemptDelay)
		attempts--
	}
	return nil, err
}

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"nineshop-be/internal/db/sqlc"
)

type Store struct {
	pool *pgxpool.Pool
	*sqlc.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		Queries: sqlc.New(pool),
	}
}
func (store *Store) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

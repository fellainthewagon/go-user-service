package postgresql

import (
	"context"
	"fmt"
	"log"
	"rest-api/internal/config"
	"rest-api/pkg/utils"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, cfg config.PostgreSQLConfig, maxAttempts int8) (pool *pgxpool.Pool, err error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, connString)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatalf("Failed to connect to postgresql: %v", err)
	}

	fmt.Println("PostgreSQL's pool has been created")
	return pool, nil
}

package lib

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/bangueco/auction-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once   sync.Once
	dbPool *pgxpool.Pool
	err    error
)

// Creates a single instance of the database connection pool that uses the pgxpool package.
// This function is thread-safe and can be called from multiple goroutines.
func InitDBConnection() (dbpool *pgxpool.Pool) {
	cfg := config.Load()

	once.Do(func() {
		dbPool, err = pgxpool.New(context.Background(), cfg.DATABASE_URL)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}

	})

	return dbPool
}

// Closes the database connection pool.
func CloseDBConnection() {

	if dbPool == nil {
		return
	}

	defer dbPool.Close()
}

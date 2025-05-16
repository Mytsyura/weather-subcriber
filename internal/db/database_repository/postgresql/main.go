package postgresql

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func getConnectionString(schemaName string) string {
	user := "mac"
	password := "postgres"
	dbAddr := "localhost:5432"
	dbDb := "weather_subscription"

	return fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s&search_path=%s",
		dbAddr,
		dbDb,
		user,
		password,
		schemaName)
}

func NewPostgresRepo(ctx context.Context, schemaName string) (*PostgresRepo, error) {
	goqu.SetDefaultPrepared(true)
	cfg, err := pgxpool.ParseConfig(getConnectionString(schemaName))
	if err != nil {
		return nil, fmt.Errorf("failed creating config: '%w'", err)
	}

	cfg.MaxConns = 500

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating pool: '%w'", err)
	}

	return &PostgresRepo{pool}, nil
}

func (p PostgresRepo) Close() {
	log.Warn().Msg("Closing db")
	p.pool.Close()
	log.Warn().Msg("Closed db")
}

func (p PostgresRepo) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

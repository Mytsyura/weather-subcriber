package infrastructure

import (
	"context"

	"github.com/rs/zerolog"

	models "weather_subscription/internal/db/models"
)

type WeatherServiceRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, email string) error
	ConfirmSubscription(ctx context.Context, token string) error
	ListActiveSubscriptions(ctx context.Context) ([]*models.Subscription, error)
	Ping(ctx context.Context) error
	Close()
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type InfraRepo interface {
	BeginTx(ctx context.Context) (Tx, error)
	Ping(ctx context.Context) error
	ExtractOrBeginTx(ctx context.Context, outerTx Tx, beginIfNotExists bool) (tx Tx, shouldCloseTx bool, err error)
	RollbackTx(ctx context.Context, tx Tx, logger zerolog.Logger)
	Close()
}

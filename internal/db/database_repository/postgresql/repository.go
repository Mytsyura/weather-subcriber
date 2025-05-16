package postgresql

import (
	"context"
	infrastructure "weather_subscription/internal/db/database_repository/infrastracture"
	models "weather_subscription/internal/db/models"
)

type postgresqlWeatherServiceRepository struct {
	repo *PostgresRepo
}

func (p postgresqlWeatherServiceRepository) Close() {
	p.repo.Close()
}

func (p postgresqlWeatherServiceRepository) Ping(ctx context.Context) error {
	return p.repo.Ping(ctx)
}

func (p postgresqlWeatherServiceRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (email, city, frequency, active, token)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := p.repo.pool.Exec(ctx, query,
		subscription.Email,
		subscription.City,
		subscription.Frequency,
		subscription.Active,
		subscription.Token,
	)

	return err
}

func (p postgresqlWeatherServiceRepository) DeleteSubscription(ctx context.Context, email string) error {
	query := `UPDATE subscriptions SET active = false WHERE email = $1 AND active = true`

	_, err := p.repo.pool.Exec(ctx, query, email)
	return err
}

func (p postgresqlWeatherServiceRepository) ConfirmSubscription(ctx context.Context, token string) error {
	query := `UPDATE subscriptions SET active = true WHERE token = $1 AND active = false`

	_, err := p.repo.pool.Exec(ctx, query, token)
	return err
}

func (p postgresqlWeatherServiceRepository) ListActiveSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE active = true`

	rows, err := p.repo.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.Email,
			&sub.City,
			&sub.Frequency,
			&sub.Token,
			&sub.Active,
			&sub.CreatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func NewWeatherServiceRepository(repo *PostgresRepo) infrastructure.WeatherServiceRepository {
	return &postgresqlWeatherServiceRepository{
		repo: repo,
	}
}

package databasehandler

import (
	"context"
	"errors"

	models "weather_subscription/internal/db/models"

	"github.com/google/uuid"
)

func CreateSubscription(ctx context.Context, email, city string, frequency models.SubscriptionFrequency) (*string, error) {

	if frequency != models.Daily && frequency != models.Hourly {
		return nil, errors.New("invalid frequency: must be 'daily' or 'hourly'")
	}
	token := uuid.New().String()
	subscription := &models.Subscription{
		Email:     email,
		City:      city,
		Frequency: frequency,
		Token:     token,
		Active:    false,
	}

	if err := dbHandler.weatherServiceRepository.CreateSubscription(ctx, subscription); err != nil {
		return nil, errors.New("subscription already exists")
	}

	return &token, nil
}

func DeleteSubscription(ctx context.Context, email string) error {

	if err := dbHandler.weatherServiceRepository.DeleteSubscription(ctx, email); err != nil {
		return errors.New("failed to delete subscription")
	}

	return nil
}

func ConfirmSubscription(ctx context.Context, token string) error {

	if err := dbHandler.weatherServiceRepository.ConfirmSubscription(ctx, token); err != nil {
		return errors.New("failed to confirm subscription")
	}

	return nil
}

func ListActiveSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	subsciptions, err := dbHandler.weatherServiceRepository.ListActiveSubscriptions(ctx)
	if err != nil {
		return nil, errors.New("failed to list active subscriptions")
	}

	return subsciptions, nil
}

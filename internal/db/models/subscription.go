package models

import (
	"time"
)

type SubscriptionFrequency string

const (
	Daily  SubscriptionFrequency = "daily"
	Hourly SubscriptionFrequency = "hourly"
)

type Subscription struct {
	ID        uint                  `json:"id"`
	Email     string                `json:"email"`
	City      string                `json:"city"`
	Frequency SubscriptionFrequency `json:"frequency"` // daily, hourly
	Token     string                `json:"token"`
	Active    bool                  `json:"active"`
	CreatedAt time.Time             `json:"created_at"`
}

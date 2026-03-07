package domain

import (
	"time"
)

type User struct {
	ID string
	Identity
	Subscription
	Devices      int
	LanguageCode string
	CurrencyCode string
	CreatedAt    time.Time
}

package domain

import (
	"time"
)

type Subscription struct {
	NodeID    string
	StartAt   time.Time
	ExpiresAt time.Time
}

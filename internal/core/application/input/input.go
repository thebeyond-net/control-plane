package input

import "time"

type NewbieConfig struct {
	Devices       int
	Bandwidth     int
	TrialDuration time.Duration
	LanguageCode  string
	CurrencyCode  string
}

type Login struct {
	ReferrerID *string
}

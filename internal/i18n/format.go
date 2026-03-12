package i18n

import "time"

func FormatRemaining(languageCode string, t time.Time) string {
	d := time.Until(t)
	if d <= 0 {
		return Get(languageCode, "SubscriptionExpired", nil, nil)
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours())
	minutes := int(d.Minutes())

	if days > 0 {
		return Get(languageCode, "DaysCount", map[string]any{"Count": days}, days)
	}
	if hours > 0 {
		return Get(languageCode, "HoursCount", map[string]any{"Count": hours}, hours)
	}

	return Get(languageCode, "MinutesCount", map[string]any{"Count": minutes}, minutes)
}

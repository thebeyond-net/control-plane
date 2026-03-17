package plans

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	sharedInput "github.com/thebeyond-net/control-plane/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

func (uc *UseCase) renderPeriodSelection(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	bandwidth int,
) error {
	currency, _ := uc.currencies.Get(user.CurrencyCode)

	text := uc.formatPlanDetails(plan, user.LanguageCode, user.CurrencyCode, bandwidth, 30)
	markup := interaction.NewReplyMarkup()

	if uc.featureFlags.IsEnabled("release.feature.bot.plans.date-picker", sharedInput.FeatureContext{
		UserID: user.ID,
	}) {
		now := time.Now()
		return uc.renderCalendar(
			ctx, msg, user,
			plan, bandwidth,
			now.Year(), now.Month(),
		)
	} else {
		for i, period := range uc.periods {
			if i%2 == 0 {
				markup.Next()
			}

			periodDiscount := int((float64(period.Discount) / 100.0) * 100)
			price, err := plan.GetPrice(user.CurrencyCode, bandwidth, period.Days, periodDiscount)
			if err != nil {
				continue
			}

			periodNameKey := i18n.Get(user.LanguageCode, period.NameKey, nil, nil)
			btnText := fmt.Sprintf("%s — %s %s", periodNameKey, formatPrice(price), currency.Symbol)
			payload := fmt.Sprintf("plan %d %d %d 0", plan.ID, bandwidth, period.Days)

			markup.AddButton(interaction.NewButton().
				Text(btnText).
				CallbackData(payload).
				Build())
		}
	}

	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData(fmt.Sprintf("plan %d %d", plan.ID, bandwidth)).
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) renderCalendar(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	bandwidth int,
	year int,
	month time.Month,
) error {
	text := i18n.Get(user.LanguageCode, "SelectPeriod", nil, nil)
	markup := interaction.NewReplyMarkup()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if year == 0 {
		year = now.Year()
		month = now.Month()
	}

	markup.Next()
	for _, dayKey := range []string{
		"WeekdayMon",
		"WeekdayTue",
		"WeekdayWed",
		"WeekdayThu",
		"WeekdayFri",
		"WeekdaySat",
		"WeekdaySun",
	} {
		btnText := i18n.Get(user.LanguageCode, dayKey, nil, nil)
		markup.AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData("ignore").
			Build())
	}

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	startOffset := (int(firstDay.Weekday()) + 6) % 7

	markup.Next()
	for i := 0; i < startOffset; i++ {
		markup.AddButton(interaction.NewButton().
			Text(" ").
			CallbackData("ignore").
			Build())
	}

	for day := 1; day <= lastDay.Day(); day++ {
		if (day+startOffset-1)%7 == 0 && day != 1 {
			markup.Next()
		}

		currentDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
		days := int(time.Until(currentDate).Hours() / 24)
		if days < 1 {
			days = 1
		}

		btnText := strconv.Itoa(day)
		payload := fmt.Sprintf("plan %d %d %d 0", plan.ID, bandwidth, days)

		if currentDate.Before(today) {
			btnText, payload = "·", "ignore"
		}

		markup.AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData(payload).
			Build())
	}

	markup.Next()

	prevM, prevY := month-1, year
	if prevM < 1 {
		prevM = 12
		prevY--
	}
	if year > now.Year() || (year == now.Year() && month > now.Month()) {
		markup.AddButton(interaction.NewButton().Text("<<").
			CallbackData(fmt.Sprintf("plan %d %d cal %d %d", plan.ID, bandwidth, prevY, int(prevM))).Build())
	} else {
		markup.AddButton(interaction.NewButton().Text(" ").CallbackData("ignore").Build())
	}

	monthName := i18n.Get(user.LanguageCode, month.String(), nil, nil)
	centerText := fmt.Sprintf("%s %d", monthName, year)
	markup.AddButton(interaction.NewButton().Text(centerText).
		CallbackData(fmt.Sprintf("plan %d %d years %d", plan.ID, bandwidth, int(month))).Build())

	nextM, nextY := month+1, year
	if nextM > 12 {
		nextM = 1
		nextY++
	}

	markup.AddButton(interaction.NewButton().Text(">>").
		CallbackData(fmt.Sprintf("plan %d %d cal %d %d", plan.ID, bandwidth, nextY, int(nextM))).Build())

	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	markup.Next().AddButton(interaction.NewButton().Text(backBtn).
		CallbackData(fmt.Sprintf("plan %d %d", plan.ID, bandwidth)).Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) renderYearSelection(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	bandwidth int,
	month time.Month,
) error {
	text := i18n.Get(user.LanguageCode, "SelectYear", nil, nil)
	markup := interaction.NewReplyMarkup()

	currentYear := time.Now().Year()

	for i := 0; i < 12; i++ {
		if i%3 == 0 {
			markup.Next()
		}
		year := currentYear + i
		markup.AddButton(interaction.NewButton().
			Text(strconv.Itoa(year)).
			CallbackData(fmt.Sprintf("plan %d %d cal %d %d", plan.ID, bandwidth, year, int(month))).
			Build())
	}

	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)
	payload := fmt.Sprintf("plan %d %d cal %d %d", plan.ID, bandwidth, currentYear, int(month))

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData(payload).
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

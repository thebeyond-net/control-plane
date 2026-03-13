package plans

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
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
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	currency, _ := uc.currencies.Get(user.CurrencyCode)

	text := uc.formatPlanDetails(plan, user.LanguageCode, user.CurrencyCode, bandwidth, 30)
	markup := interaction.NewReplyMarkup()

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

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData(fmt.Sprintf("plan %d %d", plan.ID, bandwidth)).
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

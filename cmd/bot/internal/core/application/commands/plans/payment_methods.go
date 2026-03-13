package plans

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

func (uc *UseCase) renderPaymentMethods(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	state requestState,
) error {
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := uc.formatPlanDetails(plan, user.LanguageCode, user.CurrencyCode, state.Bandwidth, state.Period)
	markup := interaction.NewReplyMarkup()

	for i, method := range uc.paymentMethods.All() {
		if i%2 == 0 {
			markup.Next()
		}

		btnText := i18n.Get(user.LanguageCode, method.Name, nil, nil)
		payload := fmt.Sprintf("plan %d %d %d %s 0", plan.ID, state.Bandwidth, state.Period, method.Code)

		markup.AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData(payload).
			IconCustomEmojiID(method.Icon).
			Build())
	}

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData(fmt.Sprintf("plan %d %d 0", plan.ID, state.Bandwidth)).
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

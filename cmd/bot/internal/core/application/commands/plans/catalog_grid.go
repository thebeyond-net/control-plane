package plans

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/pkg/keycap"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

func (uc *UseCase) renderCatalogGrid(
	ctx context.Context,
	msg input.Message,
	user domain.User,
) error {
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "SelectDevicesCount", nil, nil)
	markup := interaction.NewReplyMarkup()

	for i := range uc.plans {
		if i%2 == 0 {
			markup.Next()
		}

		btnText := keycap.Convert(fmt.Sprint(i + 1))
		payload := fmt.Sprintf("plan %d", i)

		markup.AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData(payload).
			Build())
	}

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("plan").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

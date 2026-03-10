package settings

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconLanguage = "5242469615874904546"
	IconCurrency = "5244606924745380898"
	IconBack     = "5242383235492648578"
)

type UseCase struct {
	bot ports.Bot
}

func NewSettingsUseCase(bot ports.Bot) ports.CommandHandler {
	return &UseCase{bot}
}

func (uc *UseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "Settings", nil, nil)
	markup := uc.buildReplyMarkup(user.LanguageCode)

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	languageBtn := i18n.Get(languageCode, "LanguageButton", nil, nil)
	currencyBtn := i18n.Get(languageCode, "CurrencyButton", nil, nil)
	backBtn := i18n.Get(languageCode, "BackButton", nil, nil)

	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(languageBtn).
		CallbackData("language").
		IconCustomEmojiID(IconLanguage).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(currencyBtn).
		CallbackData("currency").
		IconCustomEmojiID(IconCurrency).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("menu").
		IconCustomEmojiID(IconBack).
		Build())

	return markup.Build()
}

package commands

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconDevices = "5244702968804050765"
	IconApp     = "5242290043292260768"
)

type ConnectionUseCase struct {
	bot ports.Bot
}

func NewConnectionUseCase(bot ports.Bot) ports.CommandHandler {
	return &ConnectionUseCase{bot}
}

func (uc *ConnectionUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "SelectAction", nil, nil)
	markup := uc.buildReplyMarkup(user.LanguageCode)

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup).
		Edit(ctx, msg.ID)
}

func (uc *ConnectionUseCase) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	devicesBtn := i18n.Get(languageCode, "DevicesButton", nil, nil)
	downloadAppBtn := i18n.Get(languageCode, "DownloadAppButton", nil, nil)
	backBtn := i18n.Get(languageCode, "BackButton", nil, nil)

	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(devicesBtn).
		CallbackData("device").
		IconCustomEmojiID(IconDevices).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(downloadAppBtn).
		CallbackData("app").
		IconCustomEmojiID(IconApp).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("menu").
		IconCustomEmojiID(IconBack).
		Build())

	return markup.Build()
}

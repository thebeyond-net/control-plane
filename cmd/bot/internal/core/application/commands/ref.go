package commands

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconShare    = "5301293363406873209"
	IconWithdraw = "5303497879925595796"
)

type RefUseCase struct {
	bot         ports.Bot
	botUsername string
	supportURL  string
}

func NewRefUseCase(
	bot ports.Bot,
	botUsername,
	supportURL string,
) ports.CommandHandler {
	return &RefUseCase{bot, botUsername, supportURL}
}

func (uc *RefUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	shareBtn := i18n.Get(user.LanguageCode, "ShareButton", nil, nil)
	withdrawBtn := i18n.Get(user.LanguageCode, "WithdrawButton", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	shareURL := i18n.Get(user.LanguageCode, "ShareURL", map[string]any{
		"BotUsername": uc.botUsername,
		"UserID":      user.ID,
	}, nil)

	text := i18n.Get(user.LanguageCode, "Ref", nil, nil)
	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(shareBtn).
		URL(fmt.Sprintf("https://t.me/share/url?url=%s", shareURL)).
		IconCustomEmojiID(IconShare).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(withdrawBtn).
		URL(uc.supportURL).
		IconCustomEmojiID(IconWithdraw).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("menu").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

package settings

import (
	"context"
	"errors"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/helpers"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	sharedPorts "github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type CurrencyUseCase struct {
	bot                 ports.Bot
	userSettingsUseCase sharedPorts.UserSettingsUseCase
	currencies          map[string]domain.Currency
	layout              [][]string
}

func NewCurrencyUseCase(
	bot ports.Bot,
	userSettingsUseCase sharedPorts.UserSettingsUseCase,
	currencies map[string]domain.Currency,
	layout [][]string,
) ports.CommandHandler {
	return &CurrencyUseCase{bot, userSettingsUseCase, currencies, layout}
}

func (uc *CurrencyUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	if len(msg.Args) > 0 {
		return uc.applyChoice(ctx, msg, user)
	}
	return uc.presentOptions(ctx, msg, user)
}

func (uc *CurrencyUseCase) applyChoice(ctx context.Context, msg input.Message, user domain.User) error {
	currencyCode := msg.Args[0]
	currency, ok := uc.currencies[currencyCode]
	if !ok {
		return errors.New("currency not found")
	}

	if err := uc.userSettingsUseCase.SetCurrency(ctx, user.ID, currencyCode); err != nil {
		return err
	}

	text := i18n.Get(user.LanguageCode, "CurrencySelected", map[string]any{
		"Name": currency.GetName(),
	}, nil)

	return uc.bot.ShowNotification(ctx, msg.InteractionID, text, false)
}

func (uc *CurrencyUseCase) presentOptions(ctx context.Context, msg input.Message, user domain.User) error {
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "SelectCurrency", nil, nil)
	markup := helpers.BuildSelectionMarkup(
		uc.layout,
		uc.currencies,
		"currency",
	)

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("settings").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

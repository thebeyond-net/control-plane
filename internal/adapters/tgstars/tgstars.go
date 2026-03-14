package tgstars

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type Adapter struct {
	bot *bot.Bot
}

func New(bot *bot.Bot) ports.Invoice {
	return &Adapter{bot}
}

func (a *Adapter) NewPayment(
	ctx context.Context,
	user domain.User,
	currency string,
	devices, bandwidth, days int,
	price float64,
) (string, error) {
	title := i18n.Get(user.LanguageCode, "DaysCount", map[string]any{
		"Count": days,
	}, days)
	description := i18n.Get(user.LanguageCode, "AgreementPayment", nil, nil)
	payload := fmt.Sprintf("devices=%d;bandwidth=%d;days=%d", devices, bandwidth, days)

	payBtn := i18n.Get(user.LanguageCode, "PayButton", nil, nil)

	markup := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: payBtn, Pay: true}},
		},
	}

	_, err := a.bot.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:      user.Identity.ProviderID,
		Title:       title,
		Description: description,
		Payload:     payload,
		Currency:    currency,
		Prices: []models.LabeledPrice{
			{Label: payBtn, Amount: int(price)},
		},
		ProviderToken: "",
		ReplyMarkup:   markup,
	})

	return "", err
}

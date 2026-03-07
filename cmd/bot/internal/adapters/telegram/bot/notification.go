package bot

import (
	"context"

	"github.com/go-telegram/bot"
)

func (a *Bot) ShowNotification(
	ctx context.Context,
	interactionID,
	text string,
	alert bool,
) error {
	_, err := a.client.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: interactionID,
		Text:            text,
		ShowAlert:       alert,
	})
	return err
}

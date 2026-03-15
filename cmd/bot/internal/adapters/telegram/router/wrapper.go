package router

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/i18n"
	"go.uber.org/zap"
)

type CommandConfig struct {
	Name string
}

type Option func(*CommandConfig)

func WithName(name string) Option {
	return func(config *CommandConfig) {
		config.Name = name
	}
}

func (r *Router) CommandWrapper(commandHandler ports.CommandHandler, opts ...Option) bot.HandlerFunc {
	next := commandHandler.Execute

	cfg := &CommandConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()

		message, providerID := toInputMessage(update)
		if update.CallbackQuery != nil {
			defer b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})
		}

		var referrerID *string
		if cfg.Name == "start" && len(message.Args) > 0 {
			referrerID = &message.Args[0]
		}

		user, err := r.authUseCase.Login(ctx, "telegram", providerID, input.Login{
			ReferrerID: referrerID,
		})
		if err != nil {
			r.appLogger.Error("failed to login",
				zap.String("command", cfg.Name),
				zap.Error(err))

			languageCode := "en"
			if update.Message != nil &&
				update.Message.From != nil {
				languageCode = update.Message.From.LanguageCode
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: message.ChatID,
				Text:   i18n.Get(languageCode, "FailedToLogin", nil, nil),
			})
			return
		}

		if err := next(ctx, message, user); err != nil {
			r.appLogger.Error("failed to execute command",
				zap.String("command", cfg.Name),
				zap.Error(err))

			languageCode := "en"
			if update.Message != nil && update.Message.From != nil {
				languageCode = update.Message.From.LanguageCode
			} else if update.CallbackQuery != nil {
				languageCode = update.CallbackQuery.From.LanguageCode
			}

			text := i18n.Get(languageCode, "FailedToExecuteCommand", nil, nil)
			switch {
			case update.Message != nil:
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   text,
				})
			case update.CallbackQuery != nil:
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            text,
				})
			}
		}
	}
}

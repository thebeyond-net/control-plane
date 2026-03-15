package main

import "github.com/go-telegram/bot"

type BotContainer struct {
	Instance *bot.Bot
}

func (c *BotContainer) GetBot() *bot.Bot {
	return c.Instance
}

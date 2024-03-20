package telegram

import (
	"fmt"
	"log"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/mymmrac/telego"
)

type Bot struct {
	*telego.Bot
	*telego.Chat
}

func NewBot(config config.Config) (Bot, error) {
	bot, err := telego.NewBot(config.TelegramToken)
	if err != nil {
		return Bot{}, fmt.Errorf("new bot: %w", err)
	}

	chat, err := bot.GetChat(&telego.GetChatParams{ChatID: telego.ChatID{ID: config.ChatId}})
	if err != nil {
		return Bot{}, fmt.Errorf("new bot: %w", err)
	}

	return Bot{
		Bot:  bot,
		Chat: chat,
	}, nil
}

func (bot Bot) Send() {
	buttons := []telego.InlineKeyboardButton{
		{Text: "Go",
			URL: "http://www.example.com/"},
	}

	markup := &telego.InlineKeyboardMarkup{}

	_, err := bot.SendMessage(&telego.SendMessageParams{ChatID: bot.ChatID(), Text: "[inline URL](http://www.example.com/)", ParseMode: "markdownv2", ReplyMarkup: markup.WithInlineKeyboard(buttons)})
	if err != nil {
		log.Fatal(err)
	}
}

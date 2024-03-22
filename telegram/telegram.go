package telegram

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/AkifhanIlgaz/pau-watcher/transaction"
	"github.com/mymmrac/telego"
)

const dexScreener = "https://dexscreener.com"

type Bot struct {
	*telego.Bot
	*telego.Chat
	config.Config
	sb strings.Builder
}

func NewBot(config *config.Config) (Bot, error) {
	bot, err := telego.NewBot(config.TelegramToken)
	if err != nil {
		return Bot{}, fmt.Errorf("new bot: %w", err)
	}

	chat, err := bot.GetChat(&telego.GetChatParams{ChatID: telego.ChatID{ID: config.ChatId}})
	if err != nil {
		return Bot{}, fmt.Errorf("new bot: %w", err)
	}

	return Bot{
		Bot:    bot,
		Chat:   chat,
		Config: *config,
		sb:     strings.Builder{},
	}, nil
}

func (bot Bot) Send(tx transaction.Transaction) {
	message := bot.generateMessage(tx)
	markup := bot.generateMarkup(tx)

	_, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID:      bot.ChatID(),
		Text:        message,
		ParseMode:   "html",
		ReplyMarkup: markup,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (bot Bot) generateDexUrl(token transaction.Token) string {
	dexUrl, _ := url.Parse(dexScreener)
	dexUrl = dexUrl.JoinPath(bot.Chain.Name, token.Address)

	query := url.Values{}
	query.Add("maker", bot.SearchAddress)

	dexUrl.RawQuery = query.Encode()

	return dexUrl.String()
}

func (bot Bot) generateSwapUrl(token transaction.Token, txType string) string {
	swapUrl, _ := url.Parse(bot.Chain.Swap.Url)

	query := url.Values{}
	query.Add("chain", bot.Chain.Name)
	if txType == "IN" {
		query.Add("inputCurrency", "eth")
		query.Add("outputCurrency", token.Address)
	} else {
		query.Add("inputCurrency", token.Address)
		query.Add("outputCurrency", "eth")
	}

	return swapUrl.String() + "?" + query.Encode()
}

func (bot Bot) generateMarkup(tx transaction.Transaction) *telego.InlineKeyboardMarkup {
	markup := &telego.InlineKeyboardMarkup{}
	buttons := []telego.InlineKeyboardButton{
		{
			Text: "Trade on " + bot.Chain.Swap.Name,
			URL:  bot.generateSwapUrl(tx.Token, tx.Type),
		},
		{
			Text: "See on DexScreener",
			URL:  bot.generateDexUrl(tx.Token),
		},
	}

	return markup.WithInlineKeyboard(buttons)
}

// TODO: Did builder reset ?
// TODO: Capitalize chain name
func (bot Bot) generateMessage(tx transaction.Transaction) string {
	bot.sb.WriteString(fmt.Sprintf("<b>%v</b>\n", tx.Token.Name))
	bot.sb.WriteString(fmt.Sprintf("<b>Chain</b>: %v\n", bot.Chain.Name))
	if tx.Type == "IN" {
		bot.sb.WriteString("🟢 <i>Buy</i>")
	} else {
		bot.sb.WriteString("🔴 <i>Sell</i>")
	}

	defer bot.sb.Reset()

	return bot.sb.String()
}

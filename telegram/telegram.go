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

const uniswap = "https://app.uniswap.org/swap"
const dexScreener = "https://dexscreener.com"

type Bot struct {
	*telego.Bot
	*telego.Chat
	config.Config
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
		Bot:    bot,
		Chat:   chat,
		Config: config,
	}, nil
}

func (bot Bot) Send(tx *transaction.Transaction) {
	markup := &telego.InlineKeyboardMarkup{}
	buttons := []telego.InlineKeyboardButton{
		{Text: "Trade on UniSwap",
			// TODO: Get chain from tx
			URL: generateSwapUrl(uniswap, "base", tx.Token, tx.Type)},
		{Text: "See on DexScreener",
			URL: generateDexUrl("base", tx.Token, bot.Pau)},
	}

	var sb strings.Builder
	sb.WriteString("<i>CALL</i>\n")
	sb.WriteString(`<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>`)
	sb.WriteString("😀")
	sb.WriteString(fmt.Sprintf("<b>%v</b>\n", tx.Token.Name))
	sb.WriteString(fmt.Sprintf("<b>%v</b>\n", tx.Type))

	defer sb.Reset()

	_, err := bot.SendMessage(&telego.SendMessageParams{ChatID: bot.ChatID(), Text: sb.String(), ParseMode: "html", ReplyMarkup: markup.WithInlineKeyboard(buttons)})
	if err != nil {
		log.Fatal(err)
	}
}

func generateDexUrl(chain string, token transaction.Token, maker string) string {
	dexUrl, _ := url.Parse(dexScreener)
	dexUrl = dexUrl.JoinPath(chain, token.Address)

	query := url.Values{}
	query.Add("maker", maker)

	dexUrl.RawQuery = query.Encode()

	return dexUrl.String()
}

func generateSwapUrl(swapSite string, chain string, token transaction.Token, txType string) string {
	swapUrl, _ := url.Parse(swapSite)

	query := url.Values{}
	query.Add("chain", chain)
	if txType == "IN" {
		query.Add("inputCurrency", "eth")
		query.Add("outputCurrency", token.Address)
	} else {
		query.Add("inputCurrency", token.Address)
		query.Add("outputCurrency", "eth")
	}

	swapUrl.RawQuery = query.Encode()

	return swapUrl.String()
}

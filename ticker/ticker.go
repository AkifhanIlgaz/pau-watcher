package ticker

import (
	"fmt"
	"log"
	"time"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/AkifhanIlgaz/pau-watcher/telegram"
	"github.com/AkifhanIlgaz/pau-watcher/transaction"
)

type Ticker struct {
	bot      telegram.Bot
	parser   transaction.Parser
	interval time.Duration
	chain    string
	lastTx   transaction.Transaction
}

func NewTicker(cfg *config.Config, interval time.Duration) *Ticker {
	parser := transaction.NewParser(cfg)
	bot, err := telegram.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &Ticker{
		bot:      bot,
		parser:   parser,
		interval: interval,
		chain:    cfg.Chain,
		lastTx:   transaction.Transaction{Timestamp: time.Time{}},
	}
}

func (ticker *Ticker) Start() {
	c := time.Tick(ticker.interval)

	for next := range c {
		_ = next

		tx, err := ticker.parser.Parse()
		if err != nil {
			log.Println(err)
			continue
		}

		if tx.Timestamp.After(ticker.lastTx.Timestamp) {
			fmt.Println("Last tx changed")
			ticker.lastTx = tx
			ticker.bot.Send(tx, ticker.chain)
		}
	}
}

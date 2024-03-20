package main

import (
	"log"
	"time"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/AkifhanIlgaz/pau-watcher/telegram"
	"github.com/AkifhanIlgaz/pau-watcher/transaction"
)

/*
First Row

	document.querySelector("tbody tr:first-child")
*/

var lastTx = transaction.Transaction{
	Timestamp: time.Time{},
}

const pauUrl = "https://basescan.org/tokentxns?a="
const timeFormat = "2006-01-02 15:04:05"
const USDC = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

const TxBuy = "IN"
const TxSell = "OUT"

func main() {
	cfg, err := config.Load("./")
	if err != nil {
		log.Fatal(err)
	}

	parser := transaction.NewParser("base", cfg.Pau)

	tx, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.NewBot(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	bot.Send(tx)

	// c := time.Tick(10 * time.Second)

	// for next := range c {
	// 	_ = next

	// 	resp, err := http.Get(pauUrl + cfg.Pau)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	defer resp.Body.Close()

	// 	document, err := goquery.NewDocumentFromReader(resp.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	tx, err := parseLastTx(document)
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}

	// 	if tx.Token == USDC {
	// 		continue
	// 	}

	// 	if tx.Timestamp.After(lastTx.Timestamp) {
	// 		fmt.Println("Last tx changed")
	// 		lastTx = tx

	// 		bot.Send(tx)
	// 	}
	// }

}

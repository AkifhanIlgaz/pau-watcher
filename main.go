package main

import (
	"flag"
	"log"
	"time"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/AkifhanIlgaz/pau-watcher/ticker"
	"github.com/AkifhanIlgaz/pau-watcher/transaction"
)

var lastTx = transaction.Transaction{
	Timestamp: time.Time{},
}

const TxBuy = "IN"
const TxSell = "OUT"

func main() {
	cfg, err := config.Load("./")
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.Chain, "chain", "", "Chain")
	flag.Parse()

	if cfg.Chain == "" {
		log.Fatal("chain is not provided")
	}

	ticker := ticker.NewTicker(cfg, 10*time.Second)

	ticker.Start()

}

package main

import (
	"log"
	"time"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/AkifhanIlgaz/pau-watcher/ticker"
)

func main() {
	cfg, err := config.Load("./")
	if err != nil {
		log.Fatal(err)
	}

	ticker := ticker.NewTicker(cfg, 10*time.Second)
	ticker.Start()
}

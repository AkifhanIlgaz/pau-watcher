package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
)

/*
First Row

	document.querySelector("tbody tr:first-child")
*/
type Transaction struct {
	Timestamp time.Time
	Type      string
	Token     string
}

var lastTx = Transaction{
	Timestamp: time.Time{},
}

const pauUrl = "https://basescan.org/tokentxns?a=0x2433f77F39815849ede7959C7c43d876242cC4BC"
const timeFormat = "2006-01-02 15:04:05"
const USDC = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

const TxBuy = "IN"
const TxSell = "OUT"

func main() {
	godotenv.Load()

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	chatId, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telego.NewBot(telegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	chat, err := bot.GetChat(&telego.GetChatParams{ChatID: telego.ChatID{ID: int64(chatId)}})
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.SendMessage(&telego.SendMessageParams{ChatID: chat.ChatID(), Text: `<a href="http://www.example.com/">inline URL</a>`, ParseMode: "HTML"})
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := http.Get(pauUrl)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer resp.Body.Close()

	// document, err := goquery.NewDocumentFromReader(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c := time.Tick(2 * time.Second)

	// for next := range c {
	// 	_ = next

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

	// 		b, err := json.MarshalIndent(&tx, "", "  ")
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			continue
	// 		}

	// 		_, err = bot.SendMessage(&telego.SendMessageParams{ChatID: chat.ChatID(), Text: `<a href="http://www.example.com/">inline URL</a>`, ParseMode: "HTML"})
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 	}
	// }

}

func parseLastTx(document *goquery.Document) (Transaction, error) {
	var tx Transaction
	var error error

	document.Find("tbody tr:first-child").Each(func(i int, s *goquery.Selection) {
		timestamp, err := parseTimeStamp(s)
		if err != nil {
			error = err
			return
		}
		tx.Timestamp = timestamp

		tx.Type = parseTxType(s)

		token, err := parseToken(s)
		if err != nil {
			error = err
			return
		}
		tx.Token = token
	})

	return tx, error
}

func parseTimeStamp(tx *goquery.Selection) (time.Time, error) {
	var timeStamp time.Time
	var err error

	tx.Find(".showDate").First().Each(func(i int, s *goquery.Selection) {
		timeStamp, err = time.Parse(timeFormat, s.Text())

	})

	return timeStamp, err
}

func parseTxType(tx *goquery.Selection) string {
	var txType string

	tx.Find(".text-align").First().Each(func(i int, s *goquery.Selection) {
		txType = strings.TrimSpace(s.Text())
	})

	return txType
}

func parseToken(tx *goquery.Selection) (string, error) {
	var token string
	var error error

	tx.Find("td:last-child a").First().Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if !exists {
			error = errors.New("Token not found")
			return
		}

		u, err := url.Parse(val)
		if err != nil {
			error = err
			return
		}

		fields := strings.Split(u.Path, "/")

		token = fields[len(fields)-1]

	})

	return token, error
}

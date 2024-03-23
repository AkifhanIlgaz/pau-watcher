package transaction

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AkifhanIlgaz/pau-watcher/config"
	"github.com/PuerkitoBio/goquery"
)

var chainToScan = map[string]string{
	"base":   "https://basescan.org/tokentxns",
	"fantom": "https://ftmscan.com/tokentxns",
}

const timeFormat = "2006-01-02 15:04:05"

type Parser struct {
	client  *http.Client
	scanUrl string
}

func NewParser(cfg *config.Config) Parser {
	scanUrl, _ := url.Parse(cfg.Chain.Scan)

	values := url.Values{}
	values.Add("a", cfg.WatchAddress)

	scanUrl.RawQuery = values.Encode()

	return Parser{
		client:  &http.Client{Timeout: 20 * time.Second},
		scanUrl: scanUrl.String(),
	}
}

func (parser *Parser) Parse() (Transaction, error) {
	resp, err := parser.client.Get(parser.scanUrl)
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}

	defer resp.Body.Close()

	fmt.Println("url", resp.Request.URL)
	fmt.Println("status", resp.Status)

	document, err := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(document.Size())
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}

	tx, err := lastTransaction(document)
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}

	return tx, nil
}

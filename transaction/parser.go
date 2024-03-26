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

	req, err := http.NewRequest(http.MethodGet, parser.scanUrl, nil)
	if err != nil {
		return Transaction{}, fmt.Errorf("client: could not create request: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")

	resp, err := parser.client.Do(req)
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}
	fmt.Println(resp.Status)
	fmt.Println(resp.Request.UserAgent())

	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}

	tx, err := lastTransaction(document)
	if err != nil {
		return Transaction{}, fmt.Errorf("parse transaction: %w", err)
	}

	return tx, nil
}

package transaction

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// TODO: Add other chains
var chainToScan = map[string]string{
	"base":   "https://basescan.org/tokentxns",
	"fantom": "https://ftmscan.com/tokentxns",
}

const timeFormat = "2006-01-02 15:04:05"

type Parser struct {
	scanUrl string
}

func NewParser(chain string, searchAddress string) Parser {
	scanUrl, _ := url.Parse(chainToScan[chain])

	values := url.Values{}
	values.Add("a", searchAddress)

	scanUrl.RawQuery = values.Encode()

	return Parser{
		scanUrl: scanUrl.String(),
	}
}

func (parser *Parser) Parse() (*Transaction, error) {
	resp, err := http.Get(parser.scanUrl)
	if err != nil {
		return nil, fmt.Errorf("parse transaction: %w", err)
	}

	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse transaction: %w", err)
	}

	tx, err := lastTransaction(document)
	if err != nil {
		return nil, fmt.Errorf("parse transaction: %w", err)
	}

	return &tx, nil
}

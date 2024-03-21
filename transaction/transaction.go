package transaction

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const USDC = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

type Transaction struct {
	Timestamp time.Time
	Type      string
	Token     Token
}

type Token struct {
	Name    string
	Address string
}

func lastTransaction(document *goquery.Document) (Transaction, error) {
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

func parseToken(tx *goquery.Selection) (Token, error) {
	var token Token
	var error error

	tx.Find("td:last-child a").First().Each(func(i int, s *goquery.Selection) {
		token.Name = strings.TrimSpace(s.Text())

		val, exists := s.Attr("href")
		if !exists {
			error = errors.New("token not found")
			return
		}

		u, err := url.Parse(val)
		if err != nil {
			error = err
			return
		}

		fields := strings.Split(u.Path, "/")
		if fields[len(fields)-1] == USDC {
			error = errors.New("USDC transaction")
			return
		}

		token.Address = fields[len(fields)-1]
	})

	return token, error
}

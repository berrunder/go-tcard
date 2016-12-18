package tcard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Card structure
type Card struct {
	Sum        int64  `json:"CardSum"`
	EndDate    string `json:"EndDate"`
	LastUsed   string `json:"Time"`
	TicketType string `json:"TicketTypeDesc"`
	City       string `json:"CityName"`
	PAN        string `json:"CardPAN"`
}

// DefaultURL to fetch card data
const DefaultURL = "https://t-karta.ru/EK/Cabinet/Trip"

// Fetch card info from url. If formURL is empty, default value will be used
func Fetch(num string, formURL string) (*Card, error) {
	if formURL == "" {
		formURL = DefaultURL
	}

	cardJSON, err := fetchJSON(num, formURL)
	if err != nil {
		return nil, err
	}

	card := Card{}
	err = json.Unmarshal(cardJSON, &card)

	if err != nil {
		return nil, err
	}

	if card.PAN == "" {
		return nil, errors.New("Card data is empty")
	}

	return &card, nil
}

func fetchJSON(num string, formURL string) ([]byte, error) {
	now := time.Now()
	currDate := now.Format("02.01.2006")
	resp, err := http.PostForm(formURL, url.Values{"pan": {num}, "currDate": {currDate}})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return extractJSON(body)
}

func extractJSON(body []byte) ([]byte, error) {
	r := regexp.MustCompile(`JSON\.parse\('([^']+)'`)

	matches := r.FindSubmatch(body)

	if matches == nil {
		fmt.Println(string(body))
		return nil, errors.New("Can't find data in response body")
	}

	return matches[1], nil
}

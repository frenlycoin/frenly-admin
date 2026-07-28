package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

// API response structures
type cryptoApiResponse struct {
	Status  string         `json:"status"`
	Symbols []cryptoSymbol `json:"symbols"`
}

type cryptoSymbol struct {
	Symbol                string `json:"symbol"`
	Last                  string `json:"last"`
	DailyChangePercentage string `json:"daily_change_percentage"`
}

// CryptoPrice holds a single cryptocurrency price
type CryptoPrice struct {
	Symbol string
	Last   string
}

// PricesManager stores current crypto prices in market cap order
type PricesManager struct {
	Prices       []CryptoPrice
	publishTimer *time.Timer
}

var prices *PricesManager

func initPrices() *PricesManager {
	p := &PricesManager{Prices: make([]CryptoPrice, 0)}
	go p.start()
	return p
}

func (p *PricesManager) start() {
	// Sleep until the start of the next hour (00 minutes)
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	if now.After(next) {
		next = next.Add(time.Hour)
	}
	time.Sleep(time.Until(next))

	for {
		p.fetch()
		p.publish()
		time.Sleep(time.Hour)
	}
}

func (p *PricesManager) fetch() {
	apiURL := "https://api.freecryptoapi.com/v1/getData?symbol=" + url.QueryEscape(strings.Join(strings.Fields(CryptoPrices), " "))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		loge(fmt.Errorf("failed to create crypto prices request: %w", err))
		return
	}
	req.Header.Set("Authorization", "Bearer "+conf.FreeCryptoAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		loge(fmt.Errorf("failed to fetch crypto prices: %w", err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		loge(fmt.Errorf("failed to read crypto prices response: %w", err))
		return
	}

	var apiResp cryptoApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		loge(fmt.Errorf("failed to parse crypto prices response: %w", err))
		return
	}

	if apiResp.Status != "success" {
		loge(fmt.Errorf("crypto prices API returned status: %s", apiResp.Status))
		return
	}

	// Build a lookup from API response
	priceMap := make(map[string]string, len(apiResp.Symbols))
	for _, s := range apiResp.Symbols {
		priceMap[s.Symbol] = s.Last
	}

	// Populate in market cap order (order defined by CryptoPrices constant)
	p.Prices = make([]CryptoPrice, 0, len(strings.Fields(CryptoPrices)))
	for _, sym := range strings.Fields(CryptoPrices) {
		if last, ok := priceMap[sym]; ok {
			p.Prices = append(p.Prices, CryptoPrice{Symbol: sym, Last: last})
		}
	}

	logs("Crypto prices updated successfully")
}

func formatPrice(price string) string {
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return price
	}

	// Split integer and decimal parts
	intPart := int64(f)
	decPart := price
	if dotIdx := strings.Index(price, "."); dotIdx != -1 {
		decPart = price[dotIdx:]
	} else {
		decPart = ""
	}

	// Add commas to integer part
	intStr := strconv.FormatInt(intPart, 10)
	n := len(intStr)
	if n <= 3 {
		return intStr + decPart
	}

	var result []byte
	for i, c := range intStr {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}

	return string(result) + decPart
}

func (p *PricesManager) publish() {
	now := time.Now().UTC()
	msg := fmt.Sprintf("📊 Crypto Prices (%s UTC)\n\n", now.Format("15:04"))

	for _, cp := range p.Prices {
		msg += fmt.Sprintf("%s: $%s\n", cp.Symbol, formatPrice(cp.Last))
	}

	rec := &telebot.Chat{ID: FrenlyCrypto}
	if _, err := b.Send(rec, msg, telebot.NoPreview); err != nil {
		loge(fmt.Errorf("failed to publish crypto prices: %w", err))
	}
}

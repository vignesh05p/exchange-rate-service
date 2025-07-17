package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"exchangerate/internal/model"
)

const apiKey = "b054bd15f679b434c262790df45a62ba"

func FetchConversionRate(from, to, date string, amount float64) (float64, error) {
	var apiUrl string

	// Sanitize inputs by trimming whitespace
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	date = strings.TrimSpace(date)

	if date == "" {
		apiUrl = fmt.Sprintf("https://api.coinlayer.com/convert?access_key=%s&from=%s&to=%s&amount=%.6f",
			apiKey,
			from,
			to,
			amount)
	} else {
		// Coinlayer requires separate logic for historical rates
		apiUrl = fmt.Sprintf("https://api.coinlayer.com/%s?access_key=%s&symbols=%s&base=%s",
			date,
			apiKey,
			to,
			from)
	}

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch from coinlayer: %v", err)
	}
	defer resp.Body.Close()

	if date == "" {
		var result model.CoinlayerConvertResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, fmt.Errorf("decode error: %v", err)
		}
		if !result.Success {
			return 0, fmt.Errorf("%s", result.Error.Info)
		}
		return result.Result, nil
	} else {
		var result model.CoinlayerHistoricalResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, fmt.Errorf("decode error: %v", err)
		}
		if !result.Success {
			return 0, fmt.Errorf("%s", result.Error.Info)
		}
		rate := result.Rates[to]
		return rate * amount, nil
	}
}

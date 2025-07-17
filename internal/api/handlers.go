package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exchangerate/internal/repository"
	"exchangerate/pkg"
)

type ConvertResponse struct {
	Amount float64 `json:"amount"`
	Error  string  `json:"error,omitempty"`
}

var cache = pkg.NewRateCache(1 * time.Hour)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// Get and sanitize inputs
	from := strings.TrimSpace(r.URL.Query().Get("from"))
	to := strings.TrimSpace(r.URL.Query().Get("to"))
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	amountStr := strings.TrimSpace(r.URL.Query().Get("amount"))

	// Validate required parameters
	if from == "" || to == "" {
		writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Missing or empty 'from' or 'to' parameter"})
		return
	}

	// Validate currency codes (simple validation, could be more thorough)
	if len(from) < 3 || len(to) < 3 {
		writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Invalid currency code format"})
		return
	}

	// Handle amount parameter with default value of 1.0
	amount := 1.0
	if amountStr != "" {
		var err error
		amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Invalid 'amount' parameter: must be a valid number"})
			return
		}
		if amount <= 0 {
			writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Invalid 'amount' parameter: must be greater than 0"})
			return
		}
	}

	// Validate date format if provided
	if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Invalid date format. Use YYYY-MM-DD"})
			return
		}

		// Validate date is not in the future
		if parsedDate.After(time.Now()) {
			writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Date cannot be in the future"})
			return
		}

		// Validate date is not too old (90 days limit for historical data)
		if time.Since(parsedDate) > 90*24*time.Hour {
			writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "Date exceeds 90-day historical limit"})
			return
		}
	}

	cacheKey := from + "-" + to
	if date == "" {
		if cached, found := cache.Get(cacheKey); found {
			writeJSON(w, http.StatusOK, ConvertResponse{Amount: cached * amount})
			return
		}
	}

	convertedAmount, err := repository.FetchConversionRate(from, to, date, amount)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ConvertResponse{Error: err.Error()})
		return
	}

	if date == "" {
		cache.Set(cacheKey, convertedAmount/amount)
	}

	writeJSON(w, http.StatusOK, ConvertResponse{Amount: convertedAmount})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

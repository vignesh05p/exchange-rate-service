package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"exchangerate/internal/repository"
)

func TestConvertHandler(t *testing.T) {
	// Enable mocking for tests
	repository.EnableMocking(true)
	defer repository.EnableMocking(false)

	// Set default mock response for successful cases
	repository.SetMockResponse(85.0, nil)
	tests := []struct {
		name         string
		queryParams  string
		wantStatus   int
		wantContains string
		wantErrMsg   string
	}{
		{
			name:        "Valid conversion request",
			queryParams: "?from=USD&to=EUR&amount=100",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Missing parameters",
			queryParams: "",
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Missing or empty 'from' or 'to' parameter",
		},
		{
			name:        "Invalid amount",
			queryParams: "?from=USD&to=EUR&amount=invalid",
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Invalid 'amount' parameter: must be a valid number",
		},
		{
			name:        "Negative amount",
			queryParams: "?from=USD&to=EUR&amount=-100",
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Invalid 'amount' parameter: must be greater than 0",
		},
		{
			name:        "Invalid date format",
			queryParams: "?from=USD&to=EUR&date=invalid-date",
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Invalid date format. Use YYYY-MM-DD",
		},
		{
			name:        "Future date",
			queryParams: "?from=USD&to=EUR&date=" + time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Date cannot be in the future",
		},
		{
			name:        "Old date",
			queryParams: "?from=USD&to=EUR&date=" + time.Now().AddDate(0, -4, 0).Format("2006-01-02"),
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Date exceeds 90-day historical limit",
		},
		{
			name:        "Invalid currency code length",
			queryParams: "?from=U&to=EUR",
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "Invalid currency code format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/convert"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ConvertHandler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatus)
			}

			// Parse response
			var response ConvertResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatal(err)
			}

			// Check error message if expected
			if tt.wantErrMsg != "" {
				if response.Error != tt.wantErrMsg {
					t.Errorf("handler returned unexpected error message: got %v want %v", response.Error, tt.wantErrMsg)
				}
			}

			// For successful cases, verify we got a valid amount
			if tt.wantStatus == http.StatusOK {
				if response.Amount <= 0 {
					t.Errorf("handler returned invalid amount: %v", response.Amount)
				}
				if response.Error != "" {
					t.Errorf("handler returned unexpected error for success case: %v", response.Error)
				}
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		payload    interface{}
		wantStatus int
	}{
		{
			name:       "Success response",
			status:     http.StatusOK,
			payload:    ConvertResponse{Amount: 100.0},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error response",
			status:     http.StatusBadRequest,
			payload:    ConvertResponse{Error: "test error"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeJSON(rr, tt.status, tt.payload)

			if rr.Code != tt.wantStatus {
				t.Errorf("writeJSON() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("writeJSON() content-type = %v, want application/json", contentType)
			}

			var response ConvertResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("writeJSON() produced invalid JSON: %v", err)
			}
		})
	}
}

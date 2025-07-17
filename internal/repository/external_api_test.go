package repository

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchConversionRate(t *testing.T) {
	// Enable mocking for tests
	EnableMocking(true)
	defer EnableMocking(false)

	tests := []struct {
		name        string
		from        string
		to          string
		date        string
		amount      float64
		mockResp    float64
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:     "Valid current conversion",
			from:     "USD",
			to:       "EUR",
			date:     "",
			amount:   100.0,
			mockResp: 85.0, // Mock EUR/USD rate
		},
		{
			name:        "Invalid currency code",
			from:        "INVALID",
			to:          "USD",
			amount:      1.0,
			mockErr:     fmt.Errorf("invalid currency code"),
			wantErr:     true,
			errContains: "invalid currency code",
		},
		{
			name:     "Valid historical conversion",
			from:     "USD",
			to:       "EUR",
			date:     time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			amount:   1.0,
			mockResp: 0.85, // Mock historical EUR/USD rate
		},
		{
			name:        "Empty currency code",
			from:        "",
			to:          "USD",
			amount:      1.0,
			wantErr:     true,
			mockErr:     fmt.Errorf("invalid currency code length"),
			errContains: "invalid currency code length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set mock response for this test
			SetMockResponse(tt.mockResp, tt.mockErr)

			got, err := FetchConversionRate(tt.from, tt.to, tt.date, tt.amount)

			// Check error cases
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchConversionRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("error message '%v' should contain '%v'", err.Error(), tt.errContains)
				}
				return
			}

			// For successful cases, verify the response
			if !tt.wantErr {
				if got <= 0 {
					t.Errorf("FetchConversionRate() returned invalid rate = %v", got)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}

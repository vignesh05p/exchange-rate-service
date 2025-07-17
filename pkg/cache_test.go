package pkg

import (
	"testing"
	"time"
)

func TestRateCache_Set(t *testing.T) {
	cache := NewRateCache(1 * time.Hour)

	tests := []struct {
		name  string
		key   string
		value float64
	}{
		{
			name:  "Basic set operation",
			key:   "USD-EUR",
			value: 0.85,
		},
		{
			name:  "Zero value",
			key:   "EUR-JPY",
			value: 0.0,
		},
		{
			name:  "Large value",
			key:   "BTC-USD",
			value: 50000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Set(tt.key, tt.value)
			got, found := cache.Get(tt.key)
			if !found {
				t.Errorf("Value not found after Set for key %s", tt.key)
			}
			if got != tt.value {
				t.Errorf("Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestRateCache_Get(t *testing.T) {
	cache := NewRateCache(100 * time.Millisecond)

	tests := []struct {
		name      string
		key       string
		value     float64
		wait      time.Duration
		wantFound bool
	}{
		{
			name:      "Get existing value",
			key:       "USD-INR",
			value:     75.0,
			wait:      0,
			wantFound: true,
		},
		{
			name:      "Get after TTL expiration",
			key:       "EUR-USD",
			value:     1.2,
			wait:      150 * time.Millisecond,
			wantFound: false,
		},
		{
			name:      "Get non-existent key",
			key:       "NON-EXISTENT",
			value:     0.0,
			wait:      0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != 0.0 {
				cache.Set(tt.key, tt.value)
			}

			time.Sleep(tt.wait)

			got, found := cache.Get(tt.key)
			if found != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", found, tt.wantFound)
			}
			if tt.wantFound && got != tt.value {
				t.Errorf("Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestRateCache_Concurrency(t *testing.T) {
	cache := NewRateCache(1 * time.Hour)
	const numGoroutines = 100
	done := make(chan bool)

	// Test concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(val float64) {
			cache.Set("TEST-KEY", val)
			done <- true
		}(float64(i))
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Test concurrent reads
	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, found := cache.Get("TEST-KEY")
			if !found {
				t.Error("Value should be found")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

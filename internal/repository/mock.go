package repository

import "sync"

var (
	mockEnabled  bool
	mockResponse float64
	mockError    error
	mockLock     sync.RWMutex
)

// EnableMocking enables mock responses for testing
func EnableMocking(enabled bool) {
	mockLock.Lock()
	defer mockLock.Unlock()
	mockEnabled = enabled
}

// SetMockResponse sets the mock response for testing
func SetMockResponse(response float64, err error) {
	mockLock.Lock()
	defer mockLock.Unlock()
	mockResponse = response
	mockError = err
}

// getFetchResponse gets the mock response if mocking is enabled
func getFetchResponse() (bool, float64, error) {
	mockLock.RLock()
	defer mockLock.RUnlock()
	if mockEnabled {
		return true, mockResponse, mockError
	}
	return false, 0, nil
}

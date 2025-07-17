package model

type CoinlayerConvertResponse struct {
	Success bool    `json:"success"`
	Result  float64 `json:"result"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error,omitempty"`
}

type CoinlayerHistoricalResponse struct {
	Success bool               `json:"success"`
	Terms   string             `json:"terms"`
	Privacy string             `json:"privacy"`
	Target  string             `json:"target"`
	Base    string             `json:"base"`
	Date    string             `json:"date"`
	Rates   map[string]float64 `json:"rates"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error,omitempty"`
}

# 💱 Exchange Rate Service

A lightweight and scalable backend service built in Go to fetch real-time and historical currency exchange rates. Built as part of the GreedyGame Backend Assignment.

## 🚀 Features

- ✅ Convert currency between supported pairs
- ✅ Fetch historical exchange rates (within last 90 days)
- ✅ Real-time API integration using **Coinlayer**
- ✅ In-memory caching for latest rates (1-hour TTL)
- ✅ Input validation and error handling
- ✅ Clean Go architecture and modular design

## ⚙️ Tech Stack

- Language: **Go**
- API: [Coinlayer](https://coinlayer.com/documentation)
- Caching: In-memory with TTL
- Frameworks: Standard `net/http`
- Containerization: Docker

## 📌 Supported Currencies

- USD (US Dollar)  
- INR (Indian Rupee)  
- EUR (Euro)  
- JPY (Japanese Yen)  
- GBP (British Pound Sterling)  

## 📥 Installation & Run

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/exchange-rate-service.git
cd exchange-rate-service
```

### 2. Add API Key

Create a `.env` file (optional) or modify `external_api.go` with your Coinlayer API key:

```go
const apiKey = "b054bd15f679b434c262790df45a62ba"
```

> 💡 You can also use `os.Getenv("COINLAYER_API_KEY")` for production environments.

### 3. Run Locally

```bash
go run ./cmd
```

## 🐳 Docker (Optional)

### Build the Docker Image

```bash
docker build -t exchange-rate-service .
```

### Run the Container

```bash
docker run -p 8080:8080 exchange-rate-service
```

## 📬 API Endpoint

### `GET /convert`

Convert currency using latest or historical rates.

#### ✅ Query Parameters:

| Param    | Type   | Required | Description                               |
| -------- | ------ | -------- | ----------------------------------------- |
| `from`   | string | ✅ Yes    | Source currency (e.g., USD)               |
| `to`     | string | ✅ Yes    | Target currency (e.g., INR)               |
| `amount` | float  | ❌ No     | Amount to convert (default: 1.0)          |
| `date`   | string | ❌ No     | Historical date (YYYY-MM-DD, max 90 days) |

#### 🔁 Example Requests:

```bash
curl "http://localhost:8080/convert?from=USD&to=INR"
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100"
curl "http://localhost:8080/convert?from=EUR&to=JPY&amount=50&date=2024-06-01"
```

#### 📦 Example Response:

```json
{
  "amount": 8598.5503
}
```

## 🛡️ Error Handling

| Status | Error Type       | Example Message                          |
| ------ | ---------------- | ---------------------------------------- |
| 400    | Validation error | `"Missing 'from' or 'to' parameter"`     |
| 400    | Invalid date     | `"Invalid date format. Use YYYY-MM-DD"`  |
| 400    | Old date (> 90d) | `"Date exceeds 90-day historical limit"` |
| 500    | API error        | `"failed to fetch from coinlayer: ..."`  |

## 🧠 Assumptions

* Default `amount` is 1.0 if not provided
* Only the listed 5 currencies are expected for now
* Historical data access is limited to 90 days (as per requirement)
* Coinlayer API key is hardcoded for demo; can be moved to `.env`

## ✅ Future Improvements

* [ ] Prometheus + Grafana instrumentation
* [ ] Crypto support (BTC, ETH, USDT)
* [ ] Unit tests for service and cache layers
* [ ] Rate limit handling for Coinlayer
* [ ] Go-kit version for microservices architecture

## 📜 License

MIT License. Use freely, improve generously.

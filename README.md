# 🚍 Seattle Transit Backend (Go)

A lightweight, fast Go API server that proxies and normalizes data from the **OneBusAway Puget Sound API** for use by a custom transit dashboard.

Built with:

- Go 1.22+
- Chi router
- net/http standard library
- Clean architecture (internal/oba, internal/service, internal/httpapi)
- Full unit test coverage for service + OBA client

## ✨ Features

- Fetch real-time arrivals for any Puget Sound stop  
  `GET /api/stops/:stopID/arrivals`
- Predictive + scheduled arrival fallback logic
- Clean, frontend-friendly JSON
- Simple to extend (nearby stops, routes, vehicles)
- Fully testable with mocked OBA client + httptest.Server

## 🚀 Running the Server

### 1. Set environment variables

```
export ONEBUSAWAY_API_KEY=your_real_key
export PORT=8080
export ONEBUSAWAY_BASE_URL=https://api.pugetsound.onebusaway.org/api/where
```

### 2. Start the server
```
go run ./cmd/server
```
### 3. Test an endpoint
```
curl http://localhost:8080/api/stops/1_75403/arrivals | jq
```

## 🧪 Tests

Run all tests:
```
go test ./...
```

## 💡 Example Response
```
[
  {
    "routeId": "ROUTE_1",
    "routeShortName": "10",
    "headsign": "Downtown",
    "etaSeconds": 240,
    "arrivalTimeEpochMs": 1732905123000,
    "predicted": true
  }
]
```
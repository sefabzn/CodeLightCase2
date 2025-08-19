# Turkcell Recommendation System - Backend API

A Go-based recommendation engine for Turkcell mobile, home internet, and TV packages. This system analyzes customer household requirements and provides personalized package recommendations with bundle discounts.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL (via Supabase)
- Environment variables configured

### Running the Server

```bash
# Install dependencies
go mod tidy

# Set environment variables
export SUPABASE_URL="your-supabase-url"
export SUPABASE_ANON_KEY="your-anon-key"
export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key"
export PORT="8000"

# Run the server
go run ./cmd/server
```

The API will be available at `http://localhost:8000`

## 📊 API Endpoints

### Health Check

#### GET `/health`
Check API and database connectivity.

**Response:**
```json
{
  "status": "ok",
  "database": "connected",
  "service": "recommendation-api",
  "version": "1.0.0"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/health
```

---

### Coverage Information

#### GET `/api/coverage/{address_id}`
Get available technology coverage for a specific address.

**Parameters:**
- `address_id` (path): Address identifier (e.g., "A1001")

**Response:**
```json
{
  "address_id": "A1001",
  "city": "Istanbul",
  "district": "Kadıköy",
  "fiber": true,
  "vdsl": true,
  "fwa": false,
  "available_tech": ["fiber", "vdsl"]
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/api/coverage/A1001
```

---

### Installation Slots

#### GET `/api/install-slots/{address_id}?tech={technology}`
Get available installation time slots for a specific address and technology.

**Parameters:**
- `address_id` (path): Address identifier
- `tech` (query): Technology type ("fiber", "vdsl", "fwa") - defaults to "fiber"

**Response:**
```json
{
  "address_id": "A1001",
  "tech": "fiber",
  "slots": [
    {
      "slot_id": "SLOT_20241215_0900",
      "address_id": "A1001",
      "slot_start": "2024-12-15T09:00:00Z",
      "slot_end": "2024-12-15T11:00:00Z",
      "tech": "fiber",
      "available": true
    }
  ]
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8000/api/install-slots/A1001?tech=fiber"
```

---

### Package Recommendations

#### POST `/api/recommendation`
Get personalized package recommendations based on household requirements.

**Request Body:**
```json
{
  "user_id": 1,
  "address_id": "A1001",
  "household": [
    {
      "line_id": "LINE001",
      "expected_gb": 8.0,
      "expected_min": 450.0,
      "tv_hd_hours": 25.0
    },
    {
      "line_id": "LINE002",
      "expected_gb": 15.0,
      "expected_min": 300.0,
      "tv_hd_hours": 10.0
    }
  ],
  "prefer_tech": ["fiber", "vdsl", "fwa"]
}
```

**Response:**
```json
{
  "top3": [
    {
      "combo_label": "Mobile + Home Internet + TV Bundle",
      "items": {
        "mobile": [
          {
            "line_id": "LINE001",
            "plan": {
              "plan_id": 101,
              "plan_name": "Süper Birikim 8GB",
              "quota_gb": 8.0,
              "quota_min": 500.0,
              "monthly_price": 120.0,
              "overage_gb": 15.0,
              "overage_min": 0.5
            },
            "line_cost": 120.0,
            "overage_gb": 0.0,
            "overage_min": 0.0
          }
        ],
        "home": {
          "home_id": 201,
          "name": "Superonline Fiber 100 Mbps",
          "tech": "fiber",
          "down_mbps": 100,
          "monthly_price": 150.0,
          "install_fee": 100.0
        },
        "tv": {
          "tv_id": 301,
          "name": "TV+ Orta Paket",
          "hd_hours_included": 50.0,
          "monthly_price": 80.0
        }
      },
      "monthly_total": 315.0,
      "savings": 35.0,
      "reasoning": "Best value with full fiber coverage and bundle discounts applied",
      "discounts": {
        "line_discount": 0.0,
        "bundle_discount": 35.0,
        "total_discount": 35.0
      }
    }
  ]
}
```

**Field Descriptions:**
- `combo_label`: Human-readable package description
- `monthly_total`: Final monthly cost after all discounts
- `savings`: Total amount saved vs individual plans
- `reasoning`: Explanation of why this package was recommended
- `discounts`: Breakdown of applied discounts

**cURL Example:**
```bash
curl -X POST http://localhost:8000/api/recommendation \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "address_id": "A1001",
    "household": [
      {
        "line_id": "LINE001",
        "expected_gb": 8.0,
        "expected_min": 450.0,
        "tv_hd_hours": 25.0
      }
    ]
  }'
```

---

### Checkout

#### POST `/api/checkout`
Process package selection and create an order.

**Request Body:**
```json
{
  "user_id": 1,
  "selected_combo": {
    "combo_label": "Mobile + Home Internet Bundle",
    "items": { /* ... recommendation items ... */ },
    "monthly_total": 270.0,
    "savings": 30.0,
    "reasoning": "Optimized for your usage patterns",
    "discounts": { /* ... discount details ... */ }
  },
  "slot_id": "SLOT_20241215_0900",
  "address_id": "A1001"
}
```

**Response:**
```json
{
  "status": "ok",
  "order_id": "ORDER_20241201_ABC123"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "selected_combo": {
      "combo_label": "Mobile + Home Bundle",
      "monthly_total": 270.0,
      "savings": 30.0
    },
    "slot_id": "SLOT_20241215_0900",
    "address_id": "A1001"
  }'
```

## 🔍 Error Handling

All endpoints return structured error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      "expected_gb must be positive",
      "address_id is required"
    ]
  }
}
```

**Common Error Codes:**
- `VALIDATION_ERROR`: Invalid request data
- `NOT_FOUND`: Resource not found
- `DATABASE_ERROR`: Database connectivity issues
- `INTERNAL_ERROR`: Unexpected server errors

## 🧪 Testing

### Quick API Test
```bash
# Test health endpoint
curl http://localhost:8000/health

# Test coverage lookup
curl http://localhost:8000/api/coverage/A1001

# Test recommendation with sample data
curl -X POST http://localhost:8000/api/recommendation \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

### Sample Addresses for Testing
- `A1001`: Full coverage (Fiber, VDSL, FWA) - Istanbul, Kadıköy
- `A1002`: VDSL & FWA coverage - Ankara, Çankaya  
- `A1003`: FWA coverage only - Izmir, Konak
- `A1004`: Fiber & VDSL coverage - Bursa, Osmangazi
- `A1005`: Full coverage - Antalya, Muratpaşa

## 🏗️ Architecture

### Project Structure
```
backend/
├── cmd/server/          # Application entrypoint
├── internal/
│   ├── api/            # DTOs and request/response models
│   ├── db/             # Database interfaces and implementations
│   ├── handlers/       # HTTP route handlers
│   ├── models/         # Domain models
│   ├── services/       # Business logic
│   └── utils/          # Utilities (config, validation)
└── test_request.json   # Sample request for testing
```

### Core Components

1. **Recommendation Engine**: Analyzes household usage patterns and available services
2. **Coverage Service**: Determines available technologies by address
3. **Cost Calculator**: Computes pricing with discounts and overage charges
4. **Bundle Optimizer**: Finds optimal service combinations

### Business Rules

#### Discount Logic
- **Multi-line Discount**: 5% for 2 lines, 10% for 3+ lines (mobile only)
- **Bundle Discount**: 10% for mobile+home, 15% for mobile+home+TV
- **Technology Priority**: Fiber > VDSL > FWA (based on speed and reliability)

#### Plan Selection
- Mobile plans chosen to minimize overage while meeting usage requirements
- Home internet speed selected based on household HD viewing hours
- TV packages matched to total household HD hour requirements

## 🔧 Configuration

### Environment Variables
```bash
# Required
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key

# Optional
PORT=8000                    # Server port (default: 8000)
GIN_MODE=release            # Gin mode for production
```

### Database Schema
The API expects a Supabase/PostgreSQL database with the following tables:
- `users`: Customer information
- `coverage`: Address-based technology availability
- `mobile_plans`: Available mobile service plans
- `home_plans`: Home internet service plans
- `tv_plans`: TV service packages
- `bundling_rules`: Discount rules and configurations
- `household`: Customer household information
- `current_services`: Existing customer services
- `install_slots`: Available installation time slots

## 📈 Performance

- **Response Times**: < 300ms for recommendations
- **Concurrent Users**: Supports 100+ concurrent requests
- **Caching**: Built-in query result caching for frequent lookups
- **Database**: Connection pooling with automatic reconnection

## 🛠️ Development

### Building
```bash
go build -o recommendation-server ./cmd/server
```

### Testing
```bash
go test ./...
```

### Linting
```bash
golangci-lint run
```

---

**API Version**: 1.0.0  
**Last Updated**: December 2024  
**Contact**: Turkcell Engineering Team

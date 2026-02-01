# Postman Collection Examples
# PT XYZ Multifinance API

## Environment Variables

Create a Postman environment with these variables:
```json
{
  "base_url": "http://localhost:8080",
  "jwt_token": "your-jwt-token-here"
}
```

## 1. Health Check

### Request
```
GET {{base_url}}/health
```

### Response (200 OK)
```json
{
  "status": "healthy",
  "time": "2026-02-01"
}
```

## 2. Create Transaction

### Request
```
POST {{base_url}}/api/v1/transactions
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "consumer_id": 1,
  "tenor": 3,
  "otr": 10000000,
  "admin_fee": 100000,
  "installment_amount": 3500000,
  "interest_amount": 500000,
  "asset_name": "Motor Honda Beat"
}
```

### Success Response (201 Created)
```json
{
  "message": "Transaction created successfully",
  "transaction": {
    "id": 1,
    "contract_number": "TRX-20260201-abc12345",
    "consumer_id": 1,
    "tenor": 3,
    "otr": 10000000,
    "admin_fee": 100000,
    "installment_amount": 3500000,
    "interest_amount": 500000,
    "asset_name": "Motor Honda Beat",
    "created_at": "2026-02-01T10:00:00Z"
  }
}
```

### Error Response: Insufficient Limit (400 Bad Request)
```json
{
  "error": "Insufficient credit limit"
}
```

### Error Response: Consumer Not Found (404 Not Found)
```json
{
  "error": "Consumer not found"
}
```

## 3. Get Transaction by Contract Number

### Request
```
GET {{base_url}}/api/v1/transactions/detail?contract_number=TRX-20260201-abc12345
Authorization: Bearer {{jwt_token}}
```

### Success Response (200 OK)
```json
{
  "id": 1,
  "contract_number": "TRX-20260201-abc12345",
  "consumer_id": 1,
  "tenor": 3,
  "otr": 10000000,
  "admin_fee": 100000,
  "installment_amount": 3500000,
  "interest_amount": 500000,
  "asset_name": "Motor Honda Beat",
  "created_at": "2026-02-01T10:00:00Z"
}
```

### Error Response (404 Not Found)
```json
{
  "error": "Transaction not found"
}
```

## 4. Get Consumer Transactions

### Request
```
GET {{base_url}}/api/v1/transactions/consumer?consumer_id=1
Authorization: Bearer {{jwt_token}}
```

### Success Response (200 OK)
```json
{
  "consumer_id": 1,
  "transactions": [
    {
      "id": 2,
      "contract_number": "TRX-20260201-xyz789",
      "consumer_id": 1,
      "tenor": 3,
      "otr": 5000000,
      "admin_fee": 50000,
      "installment_amount": 1750000,
      "interest_amount": 250000,
      "asset_name": "Kulkas Samsung",
      "created_at": "2026-02-01T11:00:00Z"
    },
    {
      "id": 1,
      "contract_number": "TRX-20260201-abc12345",
      "consumer_id": 1,
      "tenor": 3,
      "otr": 10000000,
      "admin_fee": 100000,
      "installment_amount": 3500000,
      "interest_amount": 500000,
      "asset_name": "Motor Honda Beat",
      "created_at": "2026-02-01T10:00:00Z"
    }
  ]
}
```

## 5. Get Consumer Limits

### Request
```
GET {{base_url}}/api/v1/limits?consumer_id=1
Authorization: Bearer {{jwt_token}}
```

### Success Response (200 OK)
```json
{
  "consumer_id": 1,
  "limits": [
    {
      "id": 1,
      "consumer_id": 1,
      "tenor": 1,
      "limit_amount": 100000,
      "used_amount": 0,
      "version": 0,
      "updated_at": "2026-02-01T08:00:00Z"
    },
    {
      "id": 2,
      "consumer_id": 1,
      "tenor": 2,
      "limit_amount": 200000,
      "used_amount": 0,
      "version": 0,
      "updated_at": "2026-02-01T08:00:00Z"
    },
    {
      "id": 3,
      "consumer_id": 1,
      "tenor": 3,
      "limit_amount": 500000,
      "used_amount": 350000,
      "version": 2,
      "updated_at": "2026-02-01T11:00:00Z"
    },
    {
      "id": 4,
      "consumer_id": 1,
      "tenor": 6,
      "limit_amount": 700000,
      "used_amount": 0,
      "version": 0,
      "updated_at": "2026-02-01T08:00:00Z"
    }
  ]
}
```

## Error Responses

### 401 Unauthorized (Missing Token)
```json
{
  "error": "Missing authorization header"
}
```

### 401 Unauthorized (Invalid Token)
```json
{
  "error": "Invalid or expired token"
}
```

### 400 Bad Request (Invalid Input)
```json
{
  "error": "Invalid request body"
}
```

### 405 Method Not Allowed
```json
{
  "error": "Method not allowed"
}
```

### 429 Too Many Requests (Rate Limit)
```json
{
  "error": "Rate limit exceeded"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create transaction"
}
```

## Testing Concurrent Transactions

### Scenario: Two simultaneous transactions for same consumer/tenor

**Request 1:**
```
POST {{base_url}}/api/v1/transactions
{
  "consumer_id": 1,
  "tenor": 3,
  "otr": 200000,
  "admin_fee": 20000,
  "installment_amount": 75000,
  "interest_amount": 5000,
  "asset_name": "TV LED 32 inch"
}
```

**Request 2 (simultaneously):**
```
POST {{base_url}}/api/v1/transactions
{
  "consumer_id": 1,
  "tenor": 3,
  "otr": 150000,
  "admin_fee": 15000,
  "installment_amount": 60000,
  "interest_amount": 5000,
  "asset_name": "Mesin Cuci"
}
```

**Expected Behavior:**
- Both requests will be processed sequentially (due to FOR UPDATE lock)
- First request: Success (if limit available)
- Second request: Success or "Insufficient limit" (depending on remaining limit)
- No race condition, no lost updates
- Version number increases correctly

## Generating JWT Token (For Testing)

Since the app requires JWT token, you can add a helper endpoint or use this Go code:

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func main() {
    secret := []byte("your-secret-key-change-in-production")
    
    claims := jwt.MapClaims{
        "consumer_id": 1,
        "nik": "3201012345678901",
        "exp": time.Now().Add(24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(secret)
    
    fmt.Println("JWT Token:", tokenString)
}
```

Or use this curl command for testing without authentication (modify middleware temporarily):
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "consumer_id": 1,
    "tenor": 3,
    "otr": 10000000,
    "admin_fee": 100000,
    "installment_amount": 3500000,
    "interest_amount": 500000,
    "asset_name": "Motor Honda Beat"
  }'
```

# PT XYZ Multifinance - Clean Architecture API

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📋 Deskripsi Proyek

PT XYZ Multifinance adalah sistem pembiayaan untuk White Goods, Motor, dan Mobil yang dibangun dengan Clean Architecture dan mengadopsi best practices keamanan OWASP.

## 🏗️ Arsitektur Aplikasi

Aplikasi ini menggunakan **Clean Architecture** dengan struktur sebagai berikut:

```
┌─────────────────────────────────────────────────────────────────┐
│                         External Layer                          │
│  (Ecommerce, Web PT XYZ, Dealer Partner, Mobile Apps)          │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway / Load Balancer                │
│                     (HTTPS, Rate Limiting)                      │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Handler Layer (Controllers)                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Transaction Handler                                    │  │
│  │  • Consumer Handler                                       │  │
│  │  • Limit Handler                                          │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Middleware Security Layer                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  ✓ JWT Authentication (OWASP #1)                         │  │
│  │  ✓ Input Validation (OWASP #3)                           │  │
│  │  ✓ Security Headers (OWASP #4)                           │  │
│  │  ✓ Rate Limiting (OWASP #5)                              │  │
│  │  ✓ Logging & Monitoring (OWASP #9)                       │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Use Case Layer                            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Transaction Use Case (Business Logic)                 │  │
│  │    - CreateTransaction (with Concurrent Handling)        │  │
│  │    - GetTransactionByContract                            │  │
│  │    - GetConsumerTransactions                             │  │
│  │    - GetConsumerLimits                                   │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                           │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Transaction Repository                                │  │
│  │  • Consumer Repository                                   │  │
│  │  • Limit Repository (Optimistic Locking)                 │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Domain Layer                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • Consumer Entity                                        │  │
│  │  • ConsumerLimit Entity                                  │  │
│  │  • Transaction Entity                                    │  │
│  │  • Business Rules & Validations                          │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  • MySQL Database (ACID Compliant)                       │  │
│  │  • Connection Pool Management                            │  │
│  │  • Database Transactions (FOR UPDATE Lock)               │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Entity Relationship Diagram (ERD)

```
┌─────────────────────────────────────────────────────────────────┐
│                         CONSUMERS                               │
├─────────────────────────────────────────────────────────────────┤
│ PK │ id                BIGINT                                   │
│    │ nik               VARCHAR(16)  UNIQUE                      │
│    │ full_name         VARCHAR(100)                             │
│    │ legal_name        VARCHAR(100)                             │
│    │ birth_place       VARCHAR(50)                              │
│    │ birth_date        DATE                                     │
│    │ salary            DECIMAL(15,2)                            │
│    │ ktp_photo         TEXT                                     │
│    │ selfie_photo      TEXT                                     │
│    │ created_at        TIMESTAMP                                │
│    │ updated_at        TIMESTAMP                                │
└────┬────────────────────────────────────────────────────────────┘
     │
     │ 1:N
     │
     ├────────────────────────────────────────────┐
     │                                            │
     ▼                                            ▼
┌────────────────────────────────┐   ┌───────────────────────────────┐
│     CONSUMER_LIMITS             │   │       TRANSACTIONS            │
├────────────────────────────────┤   ├───────────────────────────────┤
│ PK │ id           BIGINT       │   │ PK │ id           BIGINT      │
│ FK │ consumer_id  BIGINT       │   │ FK │ consumer_id  BIGINT      │
│    │ tenor        INT          │   │    │ contract_number VARCHAR  │
│    │ limit_amount DECIMAL      │   │    │ tenor        INT         │
│    │ used_amount  DECIMAL      │   │    │ otr          DECIMAL     │
│    │ version      INT          │   │    │ admin_fee    DECIMAL     │
│    │ created_at   TIMESTAMP    │   │    │ installment_amount       │
│    │ updated_at   TIMESTAMP    │   │    │ interest_amount DECIMAL  │
└────────────────────────────────┘   │    │ asset_name   VARCHAR     │
                                      │    │ created_at   TIMESTAMP   │
                                      └───────────────────────────────┘

Relationship Notes:
- One CONSUMER can have multiple CONSUMER_LIMITS (one per tenor: 1, 2, 3, 6)
- One CONSUMER can have multiple TRANSACTIONS
- CONSUMER_LIMITS uses optimistic locking (version field) for concurrent updates
- All foreign keys use CASCADE on DELETE
```

## ✨ Fitur Utama

### 1. Clean Code Architecture
- **Domain Layer**: Entities dan business rules
- **Use Case Layer**: Application business logic
- **Repository Layer**: Data access abstraction
- **Handler Layer**: HTTP request handlers
- **Infrastructure Layer**: Database, security, external services

### 2. Concurrent Transaction Handling
- Database transactions dengan `FOR UPDATE` locks
- Optimistic locking menggunakan version field
- ACID compliance untuk data consistency
- Retry mechanism untuk concurrent updates

### 3. OWASP Security Implementation
Implementasi minimal 5 dari OWASP Top 10:

#### ✅ #1 Broken Access Control
- JWT Authentication middleware
- Token-based authorization
- Context-based user identification

#### ✅ #2 Cryptographic Failures
- Constant-time comparison untuk API keys
- Secure password hashing ready
- HTTPS enforcement headers

#### ✅ #3 Injection Prevention
- Input validation middleware
- Content-Type validation
- Request size limiting (10MB max)
- Parameterized SQL queries

#### ✅ #4 Security Misconfiguration
- Security headers middleware
  - X-Frame-Options: DENY
  - X-Content-Type-Options: nosniff
  - X-XSS-Protection
  - Strict-Transport-Security
  - Content-Security-Policy
  - Cache-Control for sensitive data

#### ✅ #5 Identification and Authentication Failures
- Rate limiting (100 requests/minute per IP)
- Token expiration (24 hours)
- Strong authentication requirements

#### ✅ #9 Security Logging & Monitoring
- Request/Response logging
- Performance monitoring
- Audit trail capability

### 4. High Availability Features
- Connection pooling (100 max connections)
- Health check endpoints
- Docker containerization
- Graceful error handling

## 🚀 Cara Menjalankan

### Prerequisites
- Go 1.23 atau lebih tinggi
- MySQL 8.0 atau lebih tinggi
- Docker & Docker Compose (optional)

### Menggunakan Docker (Recommended)

1. Clone repository
```bash
git clone <repository-url>
cd xyz-multifinance
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Jalankan dengan Docker Compose
```bash
docker-compose up -d
```

4. Aplikasi akan berjalan di `http://localhost:8080`

### Tanpa Docker

1. Setup database MySQL
```bash
mysql -u root -p < migrations/schema.sql
```

2. Copy environment file dan sesuaikan
```bash
cp .env.example .env
```

3. Install dependencies
```bash
go mod download
```

4. Jalankan aplikasi
```bash
go run main.go
```

## 📝 API Endpoints

### Health Check
```
GET /health
```

### Transaction Endpoints (Requires Authentication)

#### Create Transaction
```
POST /api/v1/transactions
Authorization: Bearer <jwt_token>
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

#### Get Transaction by Contract Number
```
GET /api/v1/transactions/detail?contract_number=TRX-20260201-abc123
Authorization: Bearer <jwt_token>
```

#### Get Consumer Transactions
```
GET /api/v1/transactions/consumer?consumer_id=1
Authorization: Bearer <jwt_token>
```

#### Get Consumer Limits
```
GET /api/v1/limits?consumer_id=1
Authorization: Bearer <jwt_token>
```

## 🧪 Testing

Jalankan unit tests:
```bash
go test ./test/... -v
```

Jalankan dengan coverage:
```bash
go test ./test/... -cover
```

## 📦 Struktur Project

```
xyz-multifinance/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies checksum
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose configuration
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── README.md                    # This file
│
├── internal/                    # Private application code
│   ├── domain/                  # Business entities
│   │   ├── consumer.go
│   │   ├── limit.go
│   │   └── transaction.go
│   │
│   ├── usecase/                 # Business logic
│   │   └── transaction_usecase.go
│   │
│   ├── repository/              # Data access layer
│   │   ├── consumer_repository.go
│   │   ├── limit_repository.go
│   │   └── transaction_repository.go
│   │
│   ├── handler/                 # HTTP handlers
│   │   └── transaction_handler.go
│   │
│   └── infrastructure/          # External services
│       ├── db/
│       │   └── mysql.go
│       └── security/
│           └── middleware.go
│
├── migrations/                  # Database migrations
│   └── schema.sql
│
└── test/                        # Test files
    └── transaction_test.go
```

## 🔒 Keamanan

1. **Authentication**: JWT-based authentication
2. **Authorization**: Role-based access control ready
3. **Input Validation**: Comprehensive validation on all inputs
4. **SQL Injection Prevention**: Parameterized queries
5. **Rate Limiting**: 100 requests per minute per IP
6. **Security Headers**: Full security headers implementation
7. **Logging**: Complete audit trail

## 📈 Scalability

1. **Horizontal Scaling**: Stateless API design
2. **Database Pooling**: Configurable connection pool
3. **Caching Ready**: Structure supports Redis integration
4. **Load Balancing Ready**: Docker Swarm/Kubernetes ready
5. **Microservices Ready**: Clean architecture allows easy separation

## 🎯 Business Requirements Compliance

✅ 99.9% availability - Achieved through:
- Health checks
- Connection pooling
- Error handling
- Docker containerization

✅ Proactive monitoring - Achieved through:
- Comprehensive logging
- Health check endpoints
- Audit logs table

✅ Fast deployment - Achieved through:
- Docker containerization
- CI/CD ready structure
- Environment-based configuration

✅ OWASP Security - Achieved through:
- 5+ OWASP Top 10 implementations
- Security middleware
- Input validation

✅ ACID Compliance - Achieved through:
- Database transactions
- FOR UPDATE locks
- Optimistic locking

## 👨‍💻 Development

### Git Flow
Project menggunakan Git Flow branching strategy:
- `main`: Production-ready code
- `develop`: Integration branch
- `feature/*`: Feature branches
- `hotfix/*`: Hotfix branches
- `release/*`: Release branches

### Coding Standards
- Follow Go conventions
- Use meaningful variable names
- Write unit tests for new features
- Document complex logic
- Keep functions small and focused

## 📄 License

This project is licensed under the MIT License.

## 👥 Contributors

- Development Team - PT XYZ Multifinance

## 📞 Support

For support, email support@xyz-multifinance.com

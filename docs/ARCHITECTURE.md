# Arsitektur Aplikasi PT XYZ Multifinance

## Overview
Arsitektur ini dirancang untuk memenuhi kebutuhan sistem pembiayaan yang scalable, maintainable, reliable, adaptable, secure, dan testable dengan availability 99.9%.

## 1. High-Level Architecture

```
┌────────────────────────────────────────────────────────────────────────────┐
│                            PRESENTATION LAYER                              │
│                                                                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Web PT XYZ │  │  E-Commerce  │  │   Dealer     │  │ Mobile App   │ │
│  │              │  │  Platform    │  │  Konvensional│  │              │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                  │                  │                  │         │
└─────────┼──────────────────┼──────────────────┼──────────────────┼─────────┘
          │                  │                  │                  │
          └──────────────────┴──────────────────┴──────────────────┘
                                    │
                                    │ HTTPS
                                    ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                          LOAD BALANCER / API GATEWAY                       │
│                                                                            │
│  • SSL/TLS Termination                                                    │
│  • Request Routing                                                         │
│  • Health Checks                                                           │
│  • Rate Limiting (Layer 7)                                                │
└────────────────────────────────┬───────────────────────────────────────────┘
                                 │
                                 ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER (Go API)                          │
│                                                                            │
│  ┌───────────────────────────────────────────────────────────────────┐   │
│  │                    SECURITY MIDDLEWARE                            │   │
│  │  • JWT Authentication (OWASP #1)                                 │   │
│  │  • Input Validation (OWASP #3)                                   │   │
│  │  • Security Headers (OWASP #4)                                   │   │
│  │  • Rate Limiting per IP (OWASP #5)                               │   │
│  │  • Request/Response Logging (OWASP #9)                           │   │
│  └───────────────────────────────────────────────────────────────────┘   │
│                                                                            │
│  ┌───────────────────────────────────────────────────────────────────┐   │
│  │                    HANDLER LAYER                                  │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │   │
│  │  │   Transaction   │  │    Consumer     │  │      Limit      │ │   │
│  │  │    Handler      │  │    Handler      │  │    Handler      │ │   │
│  │  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘ │   │
│  └───────────┼────────────────────┼────────────────────┼──────────┘   │
│              │                    │                    │                │
│  ┌───────────┼────────────────────┼────────────────────┼──────────┐   │
│  │           │    USE CASE LAYER  │                    │          │   │
│  │  ┌────────▼────────┐  ┌────────▼────────┐  ┌───────▼───────┐ │   │
│  │  │   Transaction   │  │    Consumer     │  │     Limit     │ │   │
│  │  │    Use Case     │  │    Use Case     │  │   Use Case    │ │   │
│  │  │                 │  │                 │  │               │ │   │
│  │  │ • Business      │  │ • Validation    │  │ • Calc Logic  │ │   │
│  │  │   Logic         │  │ • Rules         │  │ • Checks      │ │   │
│  │  │ • Concurrent    │  │                 │  │               │ │   │
│  │  │   Handling      │  │                 │  │               │ │   │
│  │  └────────┬────────┘  └────────┬────────┘  └───────┬───────┘ │   │
│  └───────────┼────────────────────┼────────────────────┼──────────┘   │
│              │                    │                    │                │
│  ┌───────────┼────────────────────┼────────────────────┼──────────┐   │
│  │           │  REPOSITORY LAYER  │                    │          │   │
│  │  ┌────────▼────────┐  ┌────────▼────────┐  ┌───────▼───────┐ │   │
│  │  │   Transaction   │  │    Consumer     │  │     Limit     │ │   │
│  │  │   Repository    │  │   Repository    │  │  Repository   │ │   │
│  │  │                 │  │                 │  │  (Optimistic  │ │   │
│  │  │ • CRUD Ops      │  │ • CRUD Ops      │  │   Locking)    │ │   │
│  │  │ • DB Tx         │  │ • Query         │  │               │ │   │
│  │  └────────┬────────┘  └────────┬────────┘  └───────┬───────┘ │   │
│  └───────────┼────────────────────┼────────────────────┼──────────┘   │
│              │                    │                    │                │
│  ┌───────────┼────────────────────┼────────────────────┼──────────┐   │
│  │           │    DOMAIN LAYER    │                    │          │   │
│  │  ┌────────▼────────┐  ┌────────▼────────┐  ┌───────▼───────┐ │   │
│  │  │   Transaction   │  │    Consumer     │  │ ConsumerLimit │ │   │
│  │  │     Entity      │  │     Entity      │  │    Entity     │ │   │
│  │  │                 │  │                 │  │               │ │   │
│  │  │ • Validation    │  │ • Validation    │  │ • Validation  │ │   │
│  │  │ • Business      │  │ • Rules         │  │ • Methods     │ │   │
│  │  │   Rules         │  │                 │  │               │ │   │
│  │  └─────────────────┘  └─────────────────┘  └───────────────┘ │   │
│  └───────────────────────────────────────────────────────────────┘   │
└────────────────────────────────┬───────────────────────────────────────────┘
                                 │
                                 ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                         INFRASTRUCTURE LAYER                               │
│                                                                            │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                      DATABASE LAYER                                │  │
│  │                                                                    │  │
│  │  ┌──────────────────────────────────────────────────────────────┐ │  │
│  │  │              MySQL 8.0 (InnoDB)                              │ │  │
│  │  │                                                              │ │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐   │ │  │
│  │  │  │  consumers  │  │ consumer_   │  │  transactions   │   │ │  │
│  │  │  │             │  │   limits    │  │                 │   │ │  │
│  │  │  │             │  │             │  │                 │   │ │  │
│  │  │  │ • ACID      │  │ • Version   │  │ • FOR UPDATE    │   │ │  │
│  │  │  │ • Indexes   │  │   Control   │  │   Lock          │   │ │  │
│  │  │  │             │  │ • Optimistic│  │                 │   │ │  │
│  │  │  │             │  │   Locking   │  │                 │   │ │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────────┘   │ │  │
│  │  └──────────────────────────────────────────────────────────────┘ │  │
│  │                                                                    │  │
│  │  Features:                                                         │  │
│  │  • Connection Pooling (Max 100 connections)                       │  │
│  │  • Transaction Isolation Level: READ COMMITTED                    │  │
│  │  • Row-level Locking for Concurrent Access                        │  │
│  │  • Foreign Keys with CASCADE                                      │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────────┘
```

## 2. Clean Architecture Layers

### 2.1 Domain Layer (Innermost)
**Responsibility**: Core business entities and rules

**Components**:
- `Consumer`: Entity konsumen dengan validasi NIK, nama, gaji
- `ConsumerLimit`: Entity limit kredit per tenor dengan version control
- `Transaction`: Entity transaksi dengan validasi business rules

**Key Features**:
- Pure business logic
- No external dependencies
- Validation methods
- Business rule enforcement

### 2.2 Use Case Layer
**Responsibility**: Application-specific business logic

**Components**:
- `TransactionUsecase`: 
  - CreateTransaction (with concurrent handling)
  - GetTransactionByContract
  - GetConsumerTransactions
  - GetConsumerLimits

**Key Features**:
- Orchestrates repositories
- Implements business workflows
- Handles database transactions
- Manages optimistic locking

### 2.3 Repository Layer
**Responsibility**: Data access abstraction

**Components**:
- `TransactionRepository`: CRUD operations untuk transactions
- `ConsumerRepository`: CRUD operations untuk consumers
- `LimitRepository`: CRUD dengan optimistic locking

**Key Features**:
- Interface-based design
- SQL parameterization (injection prevention)
- Transaction management
- Connection pooling

### 2.4 Handler Layer (Outermost)
**Responsibility**: HTTP request/response handling

**Components**:
- `TransactionHandler`: REST API endpoints
- Request validation
- Response formatting
- Error handling

**Key Features**:
- HTTP method validation
- Input sanitization
- JSON serialization
- Status code management

## 3. Security Architecture (OWASP)

### 3.1 Authentication & Authorization
```
Request → AuthMiddleware → JWT Validation → Context Injection → Handler
```

**Implementation**:
- JWT tokens with 24-hour expiry
- Bearer token format
- Consumer ID in token claims
- Context-based authorization

### 3.2 Input Validation
```
Request → InputValidationMiddleware → Content-Type Check → Size Limit → Handler
```

**Implementation**:
- Content-Type validation (application/json only)
- Request size limit (10MB)
- Domain-level validation
- SQL injection prevention

### 3.3 Security Headers
```
Response ← SecurityHeadersMiddleware ← Application Response
```

**Headers Applied**:
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security
- Content-Security-Policy
- Cache-Control (no-store for sensitive data)

### 3.4 Rate Limiting
```
Request → RateLimitMiddleware → IP Check → Request Count → Handler/Reject
```

**Implementation**:
- 100 requests per minute per IP
- In-memory tracking
- Sliding window algorithm
- 429 status code on limit exceeded

### 3.5 Logging & Monitoring
```
Request → LoggingMiddleware → Log Request → Process → Log Response → Response
```

**Logged Information**:
- Request method and path
- IP address
- Response status code
- Processing duration
- Error details

## 4. Concurrent Transaction Handling

### 4.1 Problem Statement
Multiple users dapat melakukan transaksi secara bersamaan terhadap limit yang sama.

### 4.2 Solution Architecture

```
Transaction 1                     Transaction 2
     │                                 │
     ├─ BEGIN TRANSACTION              ├─ BEGIN TRANSACTION
     │                                 │
     ├─ SELECT ... FOR UPDATE          │
     │  (Acquire Row Lock)             │
     │                                 ├─ SELECT ... FOR UPDATE
     ├─ Check Limit                    │  (Wait for lock...)
     │                                 │
     ├─ Update Used Amount             │
     │  (version check)                │
     │                                 │
     ├─ Insert Transaction             │
     │                                 │
     ├─ COMMIT                         │
     │  (Release Lock)                 ├─ (Lock acquired)
     │                                 ├─ Check Limit
     │                                 ├─ Update Used Amount
     │                                 ├─ Insert Transaction
     │                                 ├─ COMMIT
     ▼                                 ▼
```

### 4.3 Implementation Techniques

**1. Database Transaction**
```go
tx, err := db.Begin()
defer tx.Rollback()
// ... operations ...
tx.Commit()
```

**2. Pessimistic Locking (FOR UPDATE)**
```sql
SELECT * FROM consumer_limits 
WHERE consumer_id = ? AND tenor = ?
FOR UPDATE
```

**3. Optimistic Locking (Version Field)**
```sql
UPDATE consumer_limits 
SET used_amount = ?, version = version + 1
WHERE id = ? AND version = ?
```

**Benefits**:
- ACID compliance
- Prevents race conditions
- Handles high concurrency
- Automatic retry capability

## 5. Scalability Strategy

### 5.1 Horizontal Scaling
```
             ┌─────────────┐
             │ Load        │
             │ Balancer    │
             └──────┬──────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
    ┌───▼───┐   ┌───▼───┐   ┌───▼───┐
    │ API 1 │   │ API 2 │   │ API 3 │
    └───┬───┘   └───┬───┘   └───┬───┘
        │           │           │
        └───────────┼───────────┘
                    │
            ┌───────▼───────┐
            │   MySQL DB    │
            │  (Master)     │
            └───────┬───────┘
                    │
        ┌───────────┼───────────┐
    ┌───▼───┐   ┌───▼───┐   ┌───▼───┐
    │Replica│   │Replica│   │Replica│
    └───────┘   └───────┘   └───────┘
```

### 5.2 Database Optimization
- Connection pooling (100 max connections)
- Prepared statements
- Proper indexing
- Query optimization
- Read replicas for read-heavy operations

### 5.3 Caching Strategy (Future Enhancement)
```
Request → Check Cache → Cache Hit? → Return from Cache
                  │
                  └─ Cache Miss → Database → Update Cache → Return
```

## 6. High Availability (99.9%)

### 6.1 Target SLA
- **Availability**: 99.9% (8.76 hours downtime per year)
- **Response Time**: < 200ms (95th percentile)
- **Throughput**: 1000+ TPS

### 6.2 HA Components

**Application Level**:
- Multiple API instances
- Health check endpoints
- Graceful shutdown
- Circuit breaker pattern (future)

**Database Level**:
- Master-Slave replication
- Automatic failover
- Connection retry logic
- Backup and recovery

**Infrastructure Level**:
- Container orchestration (Kubernetes)
- Auto-scaling policies
- Monitoring and alerting
- Disaster recovery plan

## 7. Deployment Architecture

### 7.1 Container Architecture
```
Docker Container (xyz-api)
├── main (binary)
├── migrations/
│   └── schema.sql
└── Environment Variables
    ├── DB_HOST
    ├── DB_PORT
    ├── DB_USER
    ├── DB_PASSWORD
    └── JWT_SECRET
```

### 7.2 Docker Compose Setup
```yaml
Services:
  - API (Go Application)
  - MySQL (Database)
  
Networks:
  - xyz-network (Bridge)
  
Volumes:
  - mysql_data (Persistent)
```

### 7.3 CI/CD Pipeline (Recommended)
```
Code Push → Git → CI Server → Build → Test → Docker Build → Push Registry → Deploy
```

## 8. Monitoring & Observability

### 8.1 Health Checks
- `/health` endpoint
- Database connectivity check
- Service status monitoring

### 8.2 Logging
- Structured logging
- Request/response logs
- Error tracking
- Audit trail

### 8.3 Metrics (Future)
- Request rate
- Error rate
- Response time
- Database connection pool usage

## 9. Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.23 | Application development |
| Database | MySQL 8.0 | Data persistence |
| Auth | JWT | Authentication |
| Containerization | Docker | Deployment |
| Orchestration | Docker Compose | Local development |
| Testing | Go testing | Unit tests |

## 10. Benefits of This Architecture

### Maintainability
- Clear separation of concerns
- Easy to understand structure
- Modular components

### Scalability
- Horizontal scaling ready
- Stateless API design
- Database pooling

### Reliability
- ACID compliance
- Error handling
- Health checks

### Security
- Multiple OWASP implementations
- Input validation
- Authentication/Authorization

### Testability
- Interface-based design
- Dependency injection
- Unit test coverage

---

**Document Version**: 1.0  
**Date**: February 1, 2026  
**Author**: PT XYZ Multifinance Development Team

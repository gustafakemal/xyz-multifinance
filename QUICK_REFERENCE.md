# Quick Reference Guide
# PT XYZ Multifinance

## 🚀 Quick Start Commands

### Start the Application
```bash
# With Docker (Recommended)
docker-compose up -d

# Without Docker
go run main.go
```

### Run Tests
```bash
go test ./test/... -v
```

### Build Application
```bash
go build -o xyz-multifinance main.go
```

## 📡 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Create Transaction (Requires JWT)
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
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

### Get Transaction by Contract
```bash
curl "http://localhost:8080/api/v1/transactions/detail?contract_number=TRX-20260201-abc123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Consumer Transactions
```bash
curl "http://localhost:8080/api/v1/transactions/consumer?consumer_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Consumer Limits
```bash
curl "http://localhost:8080/api/v1/limits?consumer_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🗄️ Database Commands

### Import Schema
```bash
mysql -u root -p < migrations/schema.sql
```

### Connect to Database
```bash
mysql -u root -p xyz_multifinance
```

### View Consumers
```sql
SELECT * FROM consumers;
```

### View Limits
```sql
SELECT c.full_name, cl.tenor, cl.limit_amount, cl.used_amount 
FROM consumers c 
JOIN consumer_limits cl ON c.id = cl.consumer_id;
```

### View Transactions
```sql
SELECT t.contract_number, c.full_name, t.asset_name, t.otr 
FROM transactions t 
JOIN consumers c ON t.consumer_id = c.id;
```

## 🐳 Docker Commands

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
docker-compose logs -f api
```

### Rebuild Images
```bash
docker-compose up -d --build
```

### Enter Container
```bash
docker exec -it xyz-api sh
```

## 📊 Sample Data

### Consumer 1: Budi
- **ID**: 1
- **NIK**: 3201012345678901
- **Salary**: Rp 5.000.000
- **Limits**:
  - 1 month: Rp 100.000
  - 2 months: Rp 200.000
  - 3 months: Rp 500.000
  - 6 months: Rp 700.000

### Consumer 2: Annisa
- **ID**: 2
- **NIK**: 3201012345678902
- **Salary**: Rp 10.000.000
- **Limits**:
  - 1 month: Rp 1.000.000
  - 2 months: Rp 1.200.000
  - 3 months: Rp 1.500.000
  - 6 months: Rp 2.000.000

## 🔐 Security Features

### OWASP Implementations
1. **JWT Authentication** - Bearer token required
2. **Input Validation** - Content-Type & size checks
3. **Security Headers** - X-Frame-Options, CSP, etc.
4. **Rate Limiting** - 100 requests/minute per IP
5. **SQL Injection Prevention** - Parameterized queries
6. **Logging & Monitoring** - All requests logged

## 🧪 Testing Scenarios

### Test 1: Valid Transaction
```json
{
  "consumer_id": 1,
  "tenor": 3,
  "otr": 100000,
  "admin_fee": 10000,
  "installment_amount": 40000,
  "interest_amount": 10000,
  "asset_name": "TV LED"
}
```
**Expected**: 201 Created

### Test 2: Insufficient Limit
```json
{
  "consumer_id": 1,
  "tenor": 1,
  "otr": 500000,
  "admin_fee": 50000,
  "installment_amount": 550000,
  "interest_amount": 0,
  "asset_name": "Laptop"
}
```
**Expected**: 400 Bad Request - "Insufficient credit limit"

### Test 3: Invalid Tenor
```json
{
  "consumer_id": 1,
  "tenor": 5,
  "otr": 100000,
  "admin_fee": 10000,
  "installment_amount": 30000,
  "interest_amount": 10000,
  "asset_name": "TV LED"
}
```
**Expected**: 400 Bad Request - Invalid tenor

## 📁 Project Structure

```
xyz-multifinance/
├── main.go                     # Entry point
├── internal/
│   ├── domain/                 # Business entities
│   ├── usecase/                # Business logic
│   ├── repository/             # Data access
│   ├── handler/                # HTTP handlers
│   └── infrastructure/         # DB, Security
├── migrations/
│   └── schema.sql              # Database schema
├── test/
│   └── transaction_test.go     # Unit tests
└── docs/                       # Documentation
```

## 🔧 Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=xyz_user
DB_PASSWORD=xyz_password
DB_NAME=xyz_multifinance

# Application
PORT=8080

# Security
JWT_SECRET=your-secret-key
API_KEY=your-api-key
```

## 📞 Common Issues & Solutions

### Issue: Port 8080 already in use
```bash
# Find process using port
netstat -ano | findstr :8080

# Kill process (Windows)
taskkill /PID <process_id> /F

# Or change port in .env
PORT=8081
```

### Issue: Database connection failed
```bash
# Check MySQL is running
docker ps | grep mysql

# Check connection
mysql -h localhost -P 3306 -u xyz_user -p
```

### Issue: Tests failing
```bash
# Clean and rebuild
go clean -testcache
go test ./test/... -v
```

## 🎯 Performance Tips

1. **Connection Pooling**: Already configured (100 max connections)
2. **Indexes**: All critical queries indexed
3. **Caching**: Ready for Redis integration
4. **Load Balancing**: Multiple API instances supported
5. **Database Replication**: Master-slave ready

## 📚 Additional Resources

- Full Documentation: `/docs/`
- Architecture Design: `/docs/ARCHITECTURE.md`
- Database ERD: `/docs/ERD.md`
- API Examples: `/docs/API_EXAMPLES.md`
- Git Workflow: `/docs/GIT_FLOW.md`

## ✅ Pre-deployment Checklist

- [ ] Update JWT_SECRET in production
- [ ] Update API_KEY in production
- [ ] Configure database credentials
- [ ] Run all tests
- [ ] Build Docker images
- [ ] Set up database backups
- [ ] Configure monitoring
- [ ] Review security settings
- [ ] Load test the application
- [ ] Document API for consumers

---

**Quick Help**: For detailed information, see `README.md` and `docs/` folder.

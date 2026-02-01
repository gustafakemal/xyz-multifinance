# Entity Relationship Diagram (ERD)
# PT XYZ Multifinance Database Design

## 1. ERD Diagram (Notation: Crow's Foot)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│                              CONSUMERS                                   │
│                                                                          │
├──────────────────────────────────────────────────────────────────────────┤
│ PK │ id                 BIGINT          AUTO_INCREMENT                  │
│    │ nik                VARCHAR(16)     UNIQUE NOT NULL                 │
│    │ full_name          VARCHAR(100)    NOT NULL                        │
│    │ legal_name         VARCHAR(100)    NOT NULL                        │
│    │ birth_place        VARCHAR(50)     NOT NULL                        │
│    │ birth_date         DATE            NOT NULL                        │
│    │ salary             DECIMAL(15,2)   NOT NULL                        │
│    │ ktp_photo          TEXT            NULL                            │
│    │ selfie_photo       TEXT            NULL                            │
│    │ created_at         TIMESTAMP       DEFAULT CURRENT_TIMESTAMP       │
│    │ updated_at         TIMESTAMP       DEFAULT CURRENT_TIMESTAMP       │
│    │                                    ON UPDATE CURRENT_TIMESTAMP     │
└────┬─────────────────────────────────────────────────────────────────────┘
     │
     │ One consumer HAS MANY limits
     │ One consumer HAS MANY transactions
     │
     ├─────────────────────────┬──────────────────────────────────┐
     │                         │                                  │
     │ 1                       │ 1                                │
     │                         │                                  │
     │ ∞                       │ ∞                                │
     │                         │                                  │
┌────▼─────────────────┐  ┌───▼──────────────────────────────────────────┐
│                      │  │                                              │
│  CONSUMER_LIMITS     │  │           TRANSACTIONS                       │
│                      │  │                                              │
├──────────────────────┤  ├──────────────────────────────────────────────┤
│ PK │ id       BIGINT │  │ PK │ id                  BIGINT              │
│ FK │ consumer_id     │  │ FK │ consumer_id         BIGINT              │
│    │          BIGINT │  │    │ contract_number     VARCHAR(50) UNIQUE  │
│    │ tenor    INT    │  │    │ tenor               INT                 │
│    │ limit_amount    │  │    │ otr                 DECIMAL(15,2)       │
│    │       DECIMAL   │  │    │ admin_fee           DECIMAL(15,2)       │
│    │ used_amount     │  │    │ installment_amount  DECIMAL(15,2)       │
│    │       DECIMAL   │  │    │ interest_amount     DECIMAL(15,2)       │
│    │ version  INT    │  │    │ asset_name          VARCHAR(100)        │
│    │ created_at      │  │    │ created_at          TIMESTAMP           │
│    │       TIMESTAMP │  │    │                                          │
│    │ updated_at      │  └──────────────────────────────────────────────┘
│    │       TIMESTAMP │
└──────────────────────┘

Indexes:
- idx_nik on consumers(nik)
- idx_consumer_id on consumer_limits(consumer_id)
- idx_tenor on consumer_limits(tenor)
- unique_consumer_tenor on consumer_limits(consumer_id, tenor)
- idx_consumer_id on transactions(consumer_id)
- idx_contract_number on transactions(contract_number)
- idx_created_at on transactions(created_at)

Foreign Keys:
- consumer_limits.consumer_id → consumers.id (ON DELETE CASCADE)
- transactions.consumer_id → consumers.id (ON DELETE CASCADE)
```

## 2. Table Specifications

### 2.1 CONSUMERS Table

**Purpose**: Menyimpan data personal konsumen

**Columns**:

| Column Name | Data Type | Constraints | Description |
|------------|-----------|-------------|-------------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| nik | VARCHAR(16) | UNIQUE, NOT NULL | Nomor KTP (16 digit) |
| full_name | VARCHAR(100) | NOT NULL | Nama lengkap konsumen |
| legal_name | VARCHAR(100) | NOT NULL | Nama sesuai KTP |
| birth_place | VARCHAR(50) | NOT NULL | Tempat lahir |
| birth_date | DATE | NOT NULL | Tanggal lahir |
| salary | DECIMAL(15,2) | NOT NULL | Gaji bulanan |
| ktp_photo | TEXT | NULL | URL/Base64 foto KTP |
| selfie_photo | TEXT | NULL | URL/Base64 foto selfie |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE | Waktu update terakhir |

**Indexes**:
- PRIMARY KEY on `id`
- UNIQUE KEY on `nik`
- INDEX on `nik` for fast lookup
- INDEX on `created_at` for time-based queries

**Business Rules**:
- NIK harus unique (satu konsumen = satu NIK)
- NIK harus 16 digit angka
- Salary harus > 0
- Birth date harus di masa lalu

**Sample Data**:
```sql
INSERT INTO consumers VALUES
(1, '3201012345678901', 'Budi Santoso', 'Budi Santoso', 'Jakarta', 
 '1990-05-15', 5000000.00, 'ktp_budi.jpg', 'selfie_budi.jpg', NOW(), NOW()),
(2, '3201012345678902', 'Annisa Rahma', 'Annisa Rahma', 'Bandung', 
 '1992-08-20', 10000000.00, 'ktp_annisa.jpg', 'selfie_annisa.jpg', NOW(), NOW());
```

### 2.2 CONSUMER_LIMITS Table

**Purpose**: Menyimpan limit kredit konsumen untuk setiap tenor

**Columns**:

| Column Name | Data Type | Constraints | Description |
|------------|-----------|-------------|-------------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| consumer_id | BIGINT | FOREIGN KEY, NOT NULL | Reference ke consumers.id |
| tenor | INT | NOT NULL, CHECK IN (1,2,3,6) | Tenor dalam bulan |
| limit_amount | DECIMAL(15,2) | NOT NULL | Total limit kredit |
| used_amount | DECIMAL(15,2) | DEFAULT 0 | Jumlah limit yang sudah digunakan |
| version | INT | DEFAULT 0 | Version untuk optimistic locking |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE | Waktu update terakhir |

**Indexes**:
- PRIMARY KEY on `id`
- FOREIGN KEY on `consumer_id` references `consumers(id)`
- UNIQUE KEY on `(consumer_id, tenor)` - satu konsumen hanya punya 1 limit per tenor
- INDEX on `consumer_id` for fast joins
- INDEX on `tenor` for tenor-based queries

**Constraints**:
- `CHECK (tenor IN (1, 2, 3, 6))` - hanya tenor valid
- `CHECK (used_amount >= 0)` - used amount tidak boleh negatif
- `CHECK (used_amount <= limit_amount)` - used tidak boleh melebihi limit

**Business Rules**:
- Setiap konsumen memiliki 4 limit (tenor 1, 2, 3, dan 6 bulan)
- Used amount dimulai dari 0
- Used amount bertambah setiap kali ada transaksi
- Version field untuk concurrent transaction handling (optimistic locking)
- Saat update, version harus match untuk mencegah lost update

**Sample Data**:
```sql
-- Budi's limits
INSERT INTO consumer_limits VALUES
(1, 1, 1, 100000.00, 0.00, 0, NOW(), NOW()),
(2, 1, 2, 200000.00, 0.00, 0, NOW(), NOW()),
(3, 1, 3, 500000.00, 0.00, 0, NOW(), NOW()),
(4, 1, 6, 700000.00, 0.00, 0, NOW(), NOW());

-- Annisa's limits
INSERT INTO consumer_limits VALUES
(5, 2, 1, 1000000.00, 0.00, 0, NOW(), NOW()),
(6, 2, 2, 1200000.00, 0.00, 0, NOW(), NOW()),
(7, 2, 3, 1500000.00, 0.00, 0, NOW(), NOW()),
(8, 2, 6, 2000000.00, 0.00, 0, NOW(), NOW());
```

### 2.3 TRANSACTIONS Table

**Purpose**: Menyimpan semua transaksi pembiayaan

**Columns**:

| Column Name | Data Type | Constraints | Description |
|------------|-----------|-------------|-------------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| contract_number | VARCHAR(50) | UNIQUE, NOT NULL | Nomor kontrak transaksi |
| consumer_id | BIGINT | FOREIGN KEY, NOT NULL | Reference ke consumers.id |
| tenor | INT | NOT NULL | Tenor dalam bulan |
| otr | DECIMAL(15,2) | NOT NULL, CHECK > 0 | On The Road price |
| admin_fee | DECIMAL(15,2) | NOT NULL, CHECK >= 0 | Biaya administrasi |
| installment_amount | DECIMAL(15,2) | NOT NULL, CHECK > 0 | Jumlah cicilan per bulan |
| interest_amount | DECIMAL(15,2) | NOT NULL, CHECK >= 0 | Total bunga |
| asset_name | VARCHAR(100) | NOT NULL | Nama barang yang dibeli |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu transaksi |

**Indexes**:
- PRIMARY KEY on `id`
- UNIQUE KEY on `contract_number`
- FOREIGN KEY on `consumer_id` references `consumers(id)`
- INDEX on `consumer_id` for fast consumer lookup
- INDEX on `contract_number` for contract lookup
- INDEX on `created_at` for time-based queries and reporting

**Constraints**:
- `CHECK (otr > 0)` - OTR harus positif
- `CHECK (admin_fee >= 0)` - Admin fee tidak boleh negatif
- `CHECK (installment_amount > 0)` - Installment harus positif
- `CHECK (interest_amount >= 0)` - Interest tidak boleh negatif

**Business Rules**:
- Contract number dibuat otomatis dengan format: TRX-YYYYMMDD-UUID
- Total amount = OTR + admin_fee + interest_amount
- Total amount harus <= available limit konsumen untuk tenor yang dipilih
- Transaksi bersifat immutable (tidak bisa diupdate, hanya insert)

**Sample Transaction Scenario**:
```sql
-- Budi membeli Motor dengan tenor 3 bulan
INSERT INTO transactions VALUES
(1, 'TRX-20260201-abc123', 1, 3, 15000000.00, 150000.00, 
 5150000.00, 300000.00, 'Motor Honda Beat', NOW());

-- This will update consumer_limits:
-- consumer_id=1, tenor=3: used_amount increases by 15,450,000
-- (15,000,000 OTR + 150,000 admin + 300,000 interest)
```

## 3. Relationships

### 3.1 CONSUMERS ↔ CONSUMER_LIMITS
- **Type**: One-to-Many
- **Cardinality**: 1 consumer dapat memiliki 4 limits (1 untuk setiap tenor)
- **Delete Rule**: CASCADE (jika consumer dihapus, semua limits-nya ikut terhapus)
- **Business Logic**: 
  - Setiap consumer otomatis dibuatkan 4 limit records saat registrasi
  - Limit tidak dapat dihapus, hanya diupdate

### 3.2 CONSUMERS ↔ TRANSACTIONS
- **Type**: One-to-Many
- **Cardinality**: 1 consumer dapat memiliki banyak transactions
- **Delete Rule**: CASCADE (jika consumer dihapus, semua transaksinya ikut terhapus)
- **Business Logic**:
  - Transaction bersifat append-only (tidak dapat diupdate)
  - Transaction history penting untuk audit trail

### 3.3 CONSUMER_LIMITS ↔ TRANSACTIONS (Implicit)
- **Type**: Logical relationship (tidak ada FK langsung)
- **Relationship**: Transaction menggunakan limit berdasarkan consumer_id dan tenor
- **Business Logic**:
  - Sebelum membuat transaction, check available limit
  - Available limit = limit_amount - used_amount
  - Setelah transaction, update used_amount

## 4. Data Integrity

### 4.1 ACID Compliance

**Atomicity**:
```sql
BEGIN TRANSACTION;
  -- Lock the limit row
  SELECT * FROM consumer_limits WHERE ... FOR UPDATE;
  -- Check limit
  IF available_limit >= transaction_amount THEN
    UPDATE consumer_limits SET used_amount = ...;
    INSERT INTO transactions ...;
    COMMIT;
  ELSE
    ROLLBACK;
  END IF;
```

**Consistency**:
- Foreign key constraints ensure referential integrity
- CHECK constraints ensure business rules
- UNIQUE constraints prevent duplicates

**Isolation**:
- FOR UPDATE lock untuk row-level locking
- Optimistic locking dengan version field
- Transaction isolation level: READ COMMITTED

**Durability**:
- InnoDB engine dengan transaction log
- Binary logging enabled
- Regular backups

### 4.2 Optimistic Locking Example

```sql
-- Thread 1: Read version
SELECT id, used_amount, version FROM consumer_limits 
WHERE consumer_id = 1 AND tenor = 3;
-- Result: id=3, used_amount=0, version=0

-- Thread 2: Read version (at same time)
SELECT id, used_amount, version FROM consumer_limits 
WHERE consumer_id = 1 AND tenor = 3;
-- Result: id=3, used_amount=0, version=0

-- Thread 1: Update with version check
UPDATE consumer_limits 
SET used_amount = 300000, version = version + 1
WHERE id = 3 AND version = 0;
-- Success: 1 row affected, version is now 1

-- Thread 2: Update with version check
UPDATE consumer_limits 
SET used_amount = 200000, version = version + 1
WHERE id = 3 AND version = 0;
-- Failed: 0 rows affected (version mismatch)
-- Must retry transaction
```

## 5. Database Performance Optimization

### 5.1 Indexing Strategy

**Primary Indexes**:
- All PRIMARY KEY columns (automatic clustered index)

**Foreign Key Indexes**:
- consumer_limits.consumer_id
- transactions.consumer_id

**Business Logic Indexes**:
- consumers.nik (unique lookup)
- transactions.contract_number (unique lookup)
- transactions.created_at (time-based queries)
- consumer_limits(consumer_id, tenor) (composite unique)

### 5.2 Query Optimization

**Common Query 1: Get consumer with all limits**
```sql
SELECT c.*, cl.tenor, cl.limit_amount, cl.used_amount
FROM consumers c
LEFT JOIN consumer_limits cl ON c.id = cl.consumer_id
WHERE c.nik = ?;
-- Uses idx_nik and idx_consumer_id
```

**Common Query 2: Get consumer transaction history**
```sql
SELECT * FROM transactions
WHERE consumer_id = ?
ORDER BY created_at DESC
LIMIT 10;
-- Uses idx_consumer_id and idx_created_at
```

**Common Query 3: Check available limit (with lock)**
```sql
SELECT limit_amount, used_amount, version
FROM consumer_limits
WHERE consumer_id = ? AND tenor = ?
FOR UPDATE;
-- Uses unique_consumer_tenor index
```

### 5.3 Table Statistics (Expected for 1M consumers)

| Table | Rows | Avg Row Size | Total Size |
|-------|------|--------------|------------|
| consumers | 1,000,000 | 500 bytes | ~500 MB |
| consumer_limits | 4,000,000 | 100 bytes | ~400 MB |
| transactions | 10,000,000+ | 200 bytes | ~2 GB+ |

## 6. Data Migration & Seeding

### 6.1 Initial Schema Creation
```bash
mysql -u root -p < migrations/schema.sql
```

### 6.2 Sample Data (Development)
Already included in schema.sql:
- 2 consumers (Budi, Annisa)
- 8 consumer_limits (4 per consumer)
- 0 transactions (to be created via API)

### 6.3 Production Data Seeding Strategy
1. Bulk import consumers from existing system
2. Calculate and set limits based on credit scoring
3. Import historical transactions
4. Verify data integrity

## 7. Backup & Recovery Strategy

### 7.1 Backup Schedule
- **Full Backup**: Daily at 2 AM
- **Incremental Backup**: Every 6 hours
- **Binary Log**: Continuous (for point-in-time recovery)

### 7.2 Retention Policy
- Daily backups: Keep for 30 days
- Monthly backups: Keep for 1 year
- Annual backups: Keep for 7 years (compliance)

### 7.3 Recovery Procedures
```bash
# Full restore
mysql -u root -p xyz_multifinance < backup_YYYYMMDD.sql

# Point-in-time recovery
mysqlbinlog --start-datetime="YYYY-MM-DD HH:MM:SS" \
            --stop-datetime="YYYY-MM-DD HH:MM:SS" \
            binlog.000001 | mysql -u root -p
```

---

**Document Version**: 1.0  
**Date**: February 1, 2026  
**Author**: PT XYZ Multifinance Development Team

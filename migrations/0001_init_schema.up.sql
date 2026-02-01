-- PT XYZ Multifinance Database Schema
-- Clean Architecture with ACID compliance

CREATE DATABASE IF NOT EXISTS xyz_multifinance;
USE xyz_multifinance;

-- Table: consumers
-- Stores customer personal information
CREATE TABLE IF NOT EXISTS consumers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nik VARCHAR(16) UNIQUE NOT NULL COMMENT 'National ID Number (KTP)',
    full_name VARCHAR(100) NOT NULL COMMENT 'Full name of consumer',
    legal_name VARCHAR(100) NOT NULL COMMENT 'Legal name as per KTP',
    birth_place VARCHAR(50) NOT NULL COMMENT 'Place of birth',
    birth_date DATE NOT NULL COMMENT 'Date of birth',
    salary DECIMAL(15,2) NOT NULL COMMENT 'Monthly salary',
    ktp_photo TEXT COMMENT 'URL or base64 of KTP photo',
    selfie_photo TEXT COMMENT 'URL or base64 of selfie photo',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nik (nik),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: consumer_limits
-- Stores credit limits for each consumer by tenor
-- Includes optimistic locking with version field for concurrent transaction handling
CREATE TABLE IF NOT EXISTS consumer_limits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    consumer_id BIGINT NOT NULL,
    tenor INT NOT NULL COMMENT 'Tenor in months: 1, 2, 3, or 6',
    limit_amount DECIMAL(15,2) NOT NULL COMMENT 'Total credit limit',
    used_amount DECIMAL(15,2) DEFAULT 0 COMMENT 'Amount already used',
    version INT DEFAULT 0 COMMENT 'Optimistic locking version',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE,
    UNIQUE KEY unique_consumer_tenor (consumer_id, tenor),
    INDEX idx_consumer_id (consumer_id),
    INDEX idx_tenor (tenor),
    CHECK (tenor IN (1, 2, 3, 6)),
    CHECK (used_amount >= 0),
    CHECK (used_amount <= limit_amount)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: transactions
-- Stores all financing transactions
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    contract_number VARCHAR(50) UNIQUE NOT NULL COMMENT 'Unique contract number',
    consumer_id BIGINT NOT NULL,
    tenor INT NOT NULL COMMENT 'Tenor in months',
    otr DECIMAL(15,2) NOT NULL COMMENT 'On The Road price',
    admin_fee DECIMAL(15,2) NOT NULL COMMENT 'Administration fee',
    installment_amount DECIMAL(15,2) NOT NULL COMMENT 'Monthly installment amount',
    interest_amount DECIMAL(15,2) NOT NULL COMMENT 'Total interest amount',
    asset_name VARCHAR(100) NOT NULL COMMENT 'Name of financed asset',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE,
    INDEX idx_consumer_id (consumer_id),
    INDEX idx_contract_number (contract_number),
    INDEX idx_created_at (created_at),
    CHECK (otr > 0),
    CHECK (admin_fee >= 0),
    CHECK (installment_amount > 0),
    CHECK (interest_amount >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sample data insertion
-- Consumer: Budi
INSERT INTO consumers (nik, full_name, legal_name, birth_place, birth_date, salary, ktp_photo, selfie_photo)
VALUES ('3201012345678901', 'Budi Santoso', 'Budi Santoso', 'Jakarta', '1990-05-15', 5000000, 'ktp_budi.jpg', 'selfie_budi.jpg');

-- Consumer: Annisa
INSERT INTO consumers (nik, full_name, legal_name, birth_place, birth_date, salary, ktp_photo, selfie_photo)
VALUES ('3201012345678902', 'Annisa Rahma', 'Annisa Rahma', 'Bandung', '1992-08-20', 10000000, 'ktp_annisa.jpg', 'selfie_annisa.jpg');

-- Limits for Budi (consumer_id = 1)
INSERT INTO consumer_limits (consumer_id, tenor, limit_amount, used_amount, version)
VALUES 
    (1, 1, 100000, 0, 0),
    (1, 2, 200000, 0, 0),
    (1, 3, 500000, 0, 0),
    (1, 6, 700000, 0, 0);

-- Limits for Annisa (consumer_id = 2)
INSERT INTO consumer_limits (consumer_id, tenor, limit_amount, used_amount, version)
VALUES 
    (2, 1, 1000000, 0, 0),
    (2, 2, 1200000, 0, 0),
    (2, 3, 1500000, 0, 0),
    (2, 6, 2000000, 0, 0);

-- Create audit log table for monitoring (OWASP Security Logging)
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    table_name VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL COMMENT 'INSERT, UPDATE, DELETE',
    record_id BIGINT NOT NULL,
    user_id BIGINT COMMENT 'ID of user who performed action',
    old_values TEXT COMMENT 'JSON of old values',
    new_values TEXT COMMENT 'JSON of new values',
    ip_address VARCHAR(45) COMMENT 'IP address of requester',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_table_action (table_name, action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


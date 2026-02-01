-- PT XYZ Multifinance Database Schema
-- DOWN Migration (Drop all objects safely)

-- Drop tables with foreign keys first
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS consumer_limits;
DROP TABLE IF EXISTS consumers;

-- Optional: drop database (HATI-HATI)
-- Uncomment kalau memang mau sekalian hapus DB
-- DROP DATABASE IF EXISTS xyz_multifinance;

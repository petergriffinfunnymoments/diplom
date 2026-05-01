-- Migration: 003_create_transaction_logs_table (Down)
-- Drops the transaction_logs table

DROP TABLE IF EXISTS transaction_logs CASCADE;

-- Migration: 001_create_transactions_table.sql
-- Creates the main transactions table for storing payment data

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL,
    idempotency_key VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    amount_value BIGINT NOT NULL,
    amount_currency VARCHAR(3) NOT NULL,
    payment_method_type VARCHAR(32) NOT NULL,
    customer_email VARCHAR(255),
    customer_phone VARCHAR(32),
    card_mask VARCHAR(32),
    description TEXT,
    external_transaction_id VARCHAR(128),
    payment_system VARCHAR(64),
    token TEXT,
    fraud_check_result VARCHAR(32) DEFAULT 'PENDING',
    retry_count INTEGER DEFAULT 0,
    routing_path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_transactions_merchant_id ON transactions(merchant_id);
CREATE INDEX idx_transactions_idempotency_key ON transactions(idempotency_key);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_payment_system ON transactions(payment_system);
CREATE INDEX idx_transactions_fraud_check_result ON transactions(fraud_check_result);

-- Comment on table
COMMENT ON TABLE transactions IS 'Хранит информацию о всех платежных транзакциях';
COMMENT ON COLUMN transactions.id IS 'Уникальный идентификатор транзакции (UUID)';
COMMENT ON COLUMN transactions.merchant_id IS 'Идентификатор мерчанта';
COMMENT ON COLUMN transactions.idempotency_key IS 'Ключ идемпотентности для защиты от дублей';
COMMENT ON COLUMN transactions.status IS 'Текущий статус транзакции';
COMMENT ON COLUMN transactions.amount_value IS 'Сумма в минорной валюте (копейки/центы)';
COMMENT ON COLUMN transactions.amount_currency IS 'Валюта платежа (ISO 4217)';
COMMENT ON COLUMN transactions.payment_method_type IS 'Тип платежного метода (CARD, SBP, DIGITAL_RUBLE, WALLET)';
COMMENT ON COLUMN transactions.card_mask IS 'Маска карты (например, **** 1234)';
COMMENT ON COLUMN transactions.external_transaction_id IS 'ID транзакции во внешней платежной системе';
COMMENT ON COLUMN transactions.payment_system IS 'Наименование платежной системы';
COMMENT ON COLUMN transactions.token IS 'Токенизированные платежные данные';
COMMENT ON COLUMN transactions.fraud_check_result IS 'Результат проверки антифродом';
COMMENT ON COLUMN transactions.retry_count IS 'Количество попыток выполнения';
COMMENT ON COLUMN transactions.routing_path IS 'Путь маршрутизации через адаптеры';

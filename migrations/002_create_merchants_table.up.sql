-- Migration: 002_create_merchants_table.sql
-- Creates table for storing merchant information

CREATE TABLE IF NOT EXISTS merchants (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(32),
    status VARCHAR(32) DEFAULT 'ACTIVE',
    webhook_url TEXT,
    webhook_secret VARCHAR(255),
    allowed_payment_methods TEXT[], -- Array of allowed payment methods
    daily_limit BIGINT DEFAULT 0, -- 0 means no limit
    monthly_limit BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_merchants_api_key_hash ON merchants(api_key_hash);
CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchants_email ON merchants(email);

-- Comments
COMMENT ON TABLE merchants IS 'Информация о мерчантах (интернет-магазинах)';
COMMENT ON COLUMN merchants.api_key_hash IS 'Хэш API ключа для аутентификации';
COMMENT ON COLUMN merchants.status IS 'Статус мерчанта (ACTIVE, BLOCKED, SUSPENDED)';
COMMENT ON COLUMN merchants.webhook_url IS 'URL для отправки уведомлений о платежах';
COMMENT ON COLUMN merchants.allowed_payment_methods IS 'Разрешенные методы оплаты';
COMMENT ON COLUMN merchants.daily_limit IS 'Дневной лимит на платежи в копейках';
COMMENT ON COLUMN merchants.monthly_limit IS 'Месячный лимит на платежи в копейках';

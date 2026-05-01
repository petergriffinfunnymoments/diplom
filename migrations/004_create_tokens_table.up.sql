-- Migration: 004_create_tokens_table.sql
-- Creates table for storing payment tokens (tokenized card data)

CREATE TABLE IF NOT EXISTS tokens (
    id VARCHAR(64) PRIMARY KEY,
    token_value VARCHAR(255) NOT NULL UNIQUE,
    original_data_hash VARCHAR(255) NOT NULL, -- Хэш оригинальных данных для проверки
    payment_method_type VARCHAR(32) NOT NULL,
    card_last_four VARCHAR(4),
    card_exp_month VARCHAR(2),
    card_exp_year VARCHAR(4),
    card_brand VARCHAR(32),
    customer_id VARCHAR(64),
    merchant_id VARCHAR(64) NOT NULL,
    status VARCHAR(32) DEFAULT 'ACTIVE', -- ACTIVE, REVOKED, EXPIRED
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    revocation_reason TEXT
);

-- Indexes
CREATE INDEX idx_tokens_token_value ON tokens(token_value);
CREATE INDEX idx_tokens_customer_id ON tokens(customer_id);
CREATE INDEX idx_tokens_merchant_id ON tokens(merchant_id);
CREATE INDEX idx_tokens_status ON tokens(status);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);

-- Comments
COMMENT ON TABLE tokens IS 'Хранит токенизированные платежные данные';
COMMENT ON COLUMN tokens.token_value IS 'Уникальный токен, заменяющий чувствительные данные';
COMMENT ON COLUMN tokens.original_data_hash IS 'Хэш оригинальных данных для верификации';
COMMENT ON COLUMN tokens.card_last_four IS 'Последние 4 цифры карты (для отображения)';
COMMENT ON COLUMN tokens.customer_id IS 'ID клиента для повторных платежей';
COMMENT ON COLUMN tokens.status IS 'Статус токена';
COMMENT ON COLUMN tokens.expires_at IS 'Дата истечения срока действия токена';

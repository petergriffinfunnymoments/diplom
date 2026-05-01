-- Migration: 003_create_transaction_logs_table.sql
-- Creates table for storing transaction audit logs

CREATE TABLE IF NOT EXISTS transaction_logs (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64) NOT NULL UNIQUE,
    payment_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    service_name VARCHAR(128) NOT NULL,
    status VARCHAR(32),
    data JSONB,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_transaction_logs_payment_id ON transaction_logs(payment_id);
CREATE INDEX idx_transaction_logs_event_type ON transaction_logs(event_type);
CREATE INDEX idx_transaction_logs_service_name ON transaction_logs(service_name);
CREATE INDEX idx_transaction_logs_created_at ON transaction_logs(created_at);
CREATE INDEX idx_transaction_logs_data ON transaction_logs USING GIN(data);

-- Comments
COMMENT ON TABLE transaction_logs IS 'Журнал событий транзакций для аудита и отладки';
COMMENT ON COLUMN transaction_logs.event_id IS 'Уникальный идентификатор события';
COMMENT ON COLUMN transaction_logs.payment_id IS 'ID связанной транзакции';
COMMENT ON COLUMN transaction_logs.event_type IS 'Тип события (PAYMENT_CREATED, FRAUD_CHECK, etc.)';
COMMENT ON COLUMN transaction_logs.service_name IS 'Наименование сервиса, создавшего событие';
COMMENT ON COLUMN transaction_logs.data IS 'Дополнительные данные события в формате JSON';
COMMENT ON COLUMN transaction_logs.error_message IS 'Сообщение об ошибке, если событие неудачное';

package models

import "time"

// ============================================
// Входные модели (Request)
// ============================================

// CreatePaymentRequest - запрос на создание платежа
type CreatePaymentRequest struct {
	MerchantID     string           `json:"merchant_id" validate:"required"`
	IdempotencyKey string           `json:"idempotency_key" validate:"required,uuid"`
	PaymentInfo    PaymentInfoInput `json:"payment_info" validate:"required"`
}

// PaymentInfoInput - информация о платеже (входная)
type PaymentInfoInput struct {
	Amount            Amount            `json:"amount" validate:"required"`
	PaymentMethodType string            `json:"payment_method_type" validate:"required,oneof=CARD SBP DIGITAL_RUBLE WALLET"`
	CustomerData      CustomerDataInput `json:"customer_data" validate:"required"`
	Description       string            `json:"description" max:"255"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// Amount - сумма платежа
type Amount struct {
	Value    int64  `json:"value" validate:"required,min=1"` // Храним в копейках/центах
	Currency string `json:"currency" validate:"required,len=3,oneof=RUB USD EUR"`
}

// CustomerDataInput - данные клиента (входные)
type CustomerDataInput struct {
	Email                string `json:"email,omitempty" validate:"omitempty,email"`
	Phone                string `json:"phone,omitempty" validate:"omitempty,e164"`
	CardNumber           string `json:"card_number,omitempty" validate:"omitempty,credit_card"`
	CardExpDate          string `json:"card_date,omitempty" validate:"omitempty,datetime=MM/YY"`
	CardCVV              string `json:"cvv_code,omitempty" validate:"omitempty,len=3|len=4"`
	DigitalWalletID      string `json:"digital_wallet_id,omitempty"`
	DigitalRubleWalletID string `json:"digital_ruble_wallet_id,omitempty"`
}

// ============================================
// Выходные модели (Response)
// ============================================

// PaymentResponse - ответ системы о платеже
type PaymentResponse struct {
	ID                 string              `json:"id"`
	MerchantID         string              `json:"merchant_id"`
	IdempotencyKey     string              `json:"idempotency_key"`
	CurrentStatus      string              `json:"current_status"`
	PaymentInfo        PaymentInfoOutput   `json:"payment_info"`
	TransactionDetails *TransactionDetails `json:"transaction_details,omitempty"`
	Error              *ErrorDetail        `json:"error,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

// PaymentInfoOutput - информация о платеже (выходная)
type PaymentInfoOutput struct {
	Amount           Amount             `json:"amount"`
	PaymentMethodType string            `json:"payment_method_type"`
	CustomerMaskedData CustomerMaskedData `json:"customer_data"`
	Description      string             `json:"description"`
}

// CustomerMaskedData - маскированные данные клиента (без чувствительных данных)
type CustomerMaskedData struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	CardMask string `json:"card_mask,omitempty"` // Например: "**** 1234"
	WalletID string `json:"wallet_id,omitempty"`
}

// TransactionDetails - детали транзакции
type TransactionDetails struct {
	ExternalTransactionID string `json:"external_transaction_id,omitempty"`
	PaymentSystem         string `json:"payment_system"`
	Token                 string `json:"token,omitempty"`
	FraudCheckResult      string `json:"fraud_check_result"`
	RetryCount            int    `json:"retry_count"`
	RoutingPath           string `json:"routing_path,omitempty"`
}

// ErrorDetail - деталь ошибки
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ============================================
// Статусы транзакций
// ============================================

type TransactionStatus string

const (
	StatusPending       TransactionStatus = "PENDING"        // Ожидает обработки
	StatusValidating    TransactionStatus = "VALIDATING"     // На валидации
	StatusFraudCheck    TransactionStatus = "FRAUD_CHECK"    // Проверка антифродом
	StatusTokenizing    TransactionStatus = "TOKENIZING"     // Токенизация
	StatusProcessing    TransactionStatus = "PROCESSING"     // Обработка в ЭПС
	StatusSuccess       TransactionStatus = "SUCCESS"        // Успешно
	StatusFailed        TransactionStatus = "FAILED"         // Неудачно
	StatusCancelled     TransactionStatus = "CANCELLED"      // Отменено
	StatusRefunded      TransactionStatus = "REFUNDED"       // Возврат
	StatusRetry         TransactionStatus = "RETRY"          // Повторная попытка
)

// ============================================
// Типы платежных методов
// ============================================

type PaymentMethodType string

const (
	MethodCard        PaymentMethodType = "CARD"
	MethodSBP         PaymentMethodType = "SBP"
	MethodDigitalRuble PaymentMethodType = "DIGITAL_RUBLE"
	MethodWallet      PaymentMethodType = "WALLET"
)

// ============================================
// Результаты проверки антифродом
// ============================================

type FraudCheckResult string

const (
	FraudPassed   FraudCheckResult = "PASSED"
	FraudFailed   FraudCheckResult = "FAILED"
	FraudReview   FraudCheckResult = "REVIEW" // Требует ручной проверки
)

// ============================================
// События для логирования
// ============================================

type EventType string

const (
	EventPaymentCreated       EventType = "PAYMENT_CREATED"
	EventPaymentValidated     EventType = "PAYMENT_VALIDATED"
	EventFraudCheckStarted    EventType = "FRAUD_CHECK_STARTED"
	EventFraudCheckCompleted  EventType = "FRAUD_CHECK_COMPLETED"
	EventTokenizationStarted  EventType = "TOKENIZATION_STARTED"
	EventTokenizationCompleted EventType = "TOKENIZATION_COMPLETED"
	EventAdapterRequest       EventType = "ADAPTER_REQUEST"
	EventAdapterResponse      EventType = "ADAPTER_RESPONSE"
	EventPaymentCompleted     EventType = "PAYMENT_COMPLETED"
	EventPaymentFailed        EventType = "PAYMENT_FAILED"
	EventRetryAttempt         EventType = "RETRY_ATTEMPT"
	EventNotificationSent     EventType = "NOTIFICATION_SENT"
)

// LogEvent - событие для логирования
type LogEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   EventType              `json:"event_type"`
	PaymentID   string                 `json:"payment_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	ServiceName string                 `json:"service_name"`
}

// ============================================
// Модель транзакции в БД
// ============================================

type Transaction struct {
	ID                  string            `db:"id"`
	MerchantID          string            `db:"merchant_id"`
	IdempotencyKey      string            `db:"idempotency_key"`
	Status              TransactionStatus `db:"status"`
	AmountValue         int64             `db:"amount_value"`
	AmountCurrency      string            `db:"amount_currency"`
	PaymentMethodType   string            `db:"payment_method_type"`
	CustomerEmail       string            `db:"customer_email"`
	CustomerPhone       string            `db:"customer_phone"`
	CardMask            string            `db:"card_mask"`
	Description         string            `db:"description"`
	ExternalTransactionID string          `db:"external_transaction_id"`
	PaymentSystem       string            `db:"payment_system"`
	Token               string            `db:"token"`
	FraudCheckResult    FraudCheckResult  `db:"fraud_check_result"`
	RetryCount          int               `db:"retry_count"`
	RoutingPath         string            `db:"routing_path"`
	CreatedAt           time.Time         `db:"created_at"`
	UpdatedAt           time.Time         `db:"updated_at"`
}

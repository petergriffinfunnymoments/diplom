package errors

import "fmt"

// ============================================
// Типы ошибок платежного шлюза
// ============================================

type ErrorCode string

const (
	// Ошибки валидации
	ErrValidation         ErrorCode = "VALIDATION_ERROR"
	ErrInvalidAmount      ErrorCode = "INVALID_AMOUNT"
	ErrInvalidCurrency    ErrorCode = "INVALID_CURRENCY"
	ErrInvalidCard        ErrorCode = "INVALID_CARD"
	ErrInvalidEmail       ErrorCode = "INVALID_EMAIL"
	ErrInvalidPhone       ErrorCode = "INVALID_PHONE"
	ErrMissingField       ErrorCode = "MISSING_FIELD"
	
	// Ошибки аутентификации и авторизации
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrForbidden          ErrorCode = "FORBIDDEN"
	ErrInvalidMerchant    ErrorCode = "INVALID_MERCHANT"
	ErrInvalidAPIKey      ErrorCode = "INVALID_API_KEY"
	
	// Ошибки идемпотентности
	ErrDuplicateRequest   ErrorCode = "DUPLICATE_REQUEST"
	ErrIdempotencyKeyUsed ErrorCode = "IDEMPOTENCY_KEY_USED"
	
	// Ошибки антифрода
	ErrFraudDetected      ErrorCode = "FRAUD_DETECTED"
	ErrFraudReview        ErrorCode = "FRAUD_REVIEW"
	
	// Ошибки токенизации
	ErrTokenizationFailed ErrorCode = "TOKENIZATION_FAILED"
	ErrInvalidToken       ErrorCode = "INVALID_TOKEN"
	
	// Ошибки платежной системы
	ErrPaymentSystem      ErrorCode = "PAYMENT_SYSTEM_ERROR"
	ErrExternalService    ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrTimeout            ErrorCode = "TIMEOUT"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	
	// Ошибки транзакции
	ErrTransactionNotFound ErrorCode = "TRANSACTION_NOT_FOUND"
	ErrTransactionFailed   ErrorCode = "TRANSACTION_FAILED"
	ErrTransactionCancelled ErrorCode = "TRANSACTION_CANCELLED"
	ErrInvalidStatus       ErrorCode = "INVALID_STATUS"
	
	// Ошибки маршрутизации
	ErrNoRouteFound       ErrorCode = "NO_ROUTE_FOUND"
	ErrAdapterNotFound    ErrorCode = "ADAPTER_NOT_FOUND"
	
	// Системные ошибки
	ErrInternal           ErrorCode = "INTERNAL_ERROR"
	ErrDatabase           ErrorCode = "DATABASE_ERROR"
	ErrUnknown            ErrorCode = "UNKNOWN_ERROR"
)

// GatewayError - основная ошибка платежного шлюза
type GatewayError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	PaymentID  string    `json:"payment_id,omitempty"`
	Retryable  bool      `json:"retryable,omitempty"`
}

// Error реализует интерфейс error
func (e *GatewayError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError создает новую ошибку шлюза
func NewError(code ErrorCode, message string) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: message,
		Retryable: isRetryable(code),
	}
}

// NewErrorWithDetails создает ошибку с дополнительными деталями
func NewErrorWithDetails(code ErrorCode, message, details string) *GatewayError {
	return &GatewayError{
		Code:    code,
		Message: message,
		Details: details,
		Retryable: isRetryable(code),
	}
}

// NewErrorWithPaymentID создает ошибку с ID платежа
func NewErrorWithPaymentID(code ErrorCode, message, paymentID string) *GatewayError {
	return &GatewayError{
		Code:      code,
		Message:   message,
		PaymentID: paymentID,
		Retryable: isRetryable(code),
	}
}

// isRetryable определяет, можно ли повторить операцию при данной ошибке
func isRetryable(code ErrorCode) bool {
	retryableCodes := map[ErrorCode]bool{
		ErrTimeout:            true,
		ErrServiceUnavailable: true,
		ErrExternalService:    true,
		ErrPaymentSystem:      true,
	}
	return retryableCodes[code]
}

// ============================================
// Конструкторы распространенных ошибок
// ============================================

// Validation errors
func ErrValidationFailed(field, reason string) *GatewayError {
	return NewErrorWithDetails(ErrValidation, fmt.Sprintf("Validation failed for field '%s'", field), reason)
}

func ErrInvalidAmountValue() *GatewayError {
	return NewError(ErrInvalidAmount, "Amount must be greater than zero")
}

func ErrInvalidCurrencyCode(currency string) *GatewayError {
	return NewErrorWithDetails(ErrInvalidCurrency, "Invalid currency code", currency)
}

func ErrInvalidCardNumber() *GatewayError {
	return NewError(ErrInvalidCard, "Invalid card number")
}

func ErrInvalidCVV() *GatewayError {
	return NewError(ErrInvalidCard, "Invalid CVV code")
}

func ErrInvalidExpiryDate() *GatewayError {
	return NewError(ErrInvalidCard, "Invalid card expiry date")
}

// Auth errors
func ErrUnauthorizedAccess() *GatewayError {
	return NewError(ErrUnauthorized, "Unauthorized access. Invalid or missing API key")
}

func ErrMerchantNotFound(merchantID string) *GatewayError {
	return NewErrorWithDetails(ErrInvalidMerchant, "Merchant not found", merchantID)
}

// Idempotency errors
func ErrDuplicateIdempotencyKey(key string) *GatewayError {
	return NewErrorWithDetails(ErrIdempotencyKeyUsed, "This idempotency key was already used", key)
}

// Fraud errors
func ErrFraudSuspicion(reason string) *GatewayError {
	return NewErrorWithDetails(ErrFraudDetected, "Transaction flagged as potentially fraudulent", reason)
}

func ErrFraudRequiresReview() *GatewayError {
	return NewError(ErrFraudReview, "Transaction requires manual review")
}

// Transaction errors
func ErrTransactionNotFoundByID(id string) *GatewayError {
	return NewErrorWithDetails(ErrTransactionNotFound, "Transaction not found", id)
}

func ErrInvalidTransactionStatus(current, expected string) *GatewayError {
	return NewErrorWithDetails(ErrInvalidStatus, 
		fmt.Sprintf("Invalid transaction status: current=%s, expected=%s", current, expected), "")
}

// System errors
func ErrInternalServer(message string) *GatewayError {
	return NewErrorWithDetails(ErrInternal, "Internal server error", message)
}

func ErrDatabaseOperation(operation, details string) *GatewayError {
	return NewErrorWithDetails(ErrDatabase, fmt.Sprintf("Database operation failed: %s", operation), details)
}

func ErrServiceTimeout(serviceName string) *GatewayError {
	return NewErrorWithDetails(ErrTimeout, fmt.Sprintf("Service timeout: %s", serviceName), "")
}

// Adapter errors
func ErrAdapterNotfound(adapterID string) *GatewayError {
	return NewErrorWithDetails(ErrAdapterNotFound, "Payment adapter not found", adapterID)
}

func ErrNoRouteForPayment(method, currency string) *GatewayError {
	return NewErrorWithDetails(ErrNoRouteFound, 
		fmt.Sprintf("No route found for payment method=%s, currency=%s", method, currency), "")
}

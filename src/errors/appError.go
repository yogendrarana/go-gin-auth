package custom_errors

// custom error type
type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// implement the error interface
func (e *AppError) Error() string {
	return e.Message
}

// create a new AppError
func NewAppError(message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

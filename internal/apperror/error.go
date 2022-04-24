package apperror

import (
	"encoding/json"
)

var NotFoundError = NewAppError(nil, "Not Found", "", "0009")

type AppError struct {
	Err        error  `json:"-"`
	Message    string `json:"message"`
	DevMessage string `json:"dev_message"`
	Code       string `json:"code"`
}

func NewAppError(err error, message, devMessage, code string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		DevMessage: devMessage,
		Code:       code,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) MarshalError() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func systemError(err error) *AppError {
	return NewAppError(err, "Internal server error", err.Error(), "0000")
}

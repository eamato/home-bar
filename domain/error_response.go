package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const (
	ReasonDBError       Reason = "db_error"
	ReasonUserDuplicate Reason = "user_duplicate"
	ReasonServerError   Reason = "server_error"
	ReasonNotFound      Reason = "not_found"
	ReasonUnauthorized  Reason = "unauthorized"
)

type Reason string

type UsecaseError struct {
	Error  error
	Reason Reason
}

func NewUsecaseError(error error, reason Reason) *UsecaseError {
	return &UsecaseError{Error: error, Reason: reason}
}

func (ue *UsecaseError) ParseUsecaseErrorToRest() (int, ErrorResponse) {
	errorResponse := GetErrorResponse(ue.Error)
	switch ue.Reason {
	case ReasonServerError, ReasonDBError:
		return http.StatusInternalServerError, errorResponse
	case ReasonUserDuplicate:
		return http.StatusConflict, errorResponse
	case ReasonNotFound:
		return http.StatusNotFound, errorResponse
	case ReasonUnauthorized:
		return http.StatusUnauthorized, errorResponse
	}

	return http.StatusInternalServerError, errorResponse
}

type CustomError struct {
	message string
}

func NewCustomError(message string) *CustomError {
	return &CustomError{
		message: message,
	}
}

func (c *CustomError) Error() string {
	return c.message
}

type ErrorResponse struct {
	Errors []ErrorStruct `json:"errors"`
}

type ErrorStruct struct {
	Field   *string `json:"field"`
	Message *string `json:"message"`
}

func GetErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Errors: getErrorResponses(err),
	}
}

func getErrorResponses(err error) []ErrorStruct {
	if err == nil {
		msg := "Error is nil"
		return []ErrorStruct{
			{Message: &msg},
		}
	}

	var ve validator.ValidationErrors
	ce := &CustomError{}

	if errors.As(err, &ve) {
		out := make([]ErrorStruct, len(ve))
		for i, fe := range ve {
			field := fe.Field()
			msg := getErrorMsg(fe)
			out[i] = ErrorStruct{&field, &msg}
		}

		return out
	} else if errors.As(err, &ce) {
		msg := ce.Error()
		return []ErrorStruct{
			{Message: &msg},
		}
	}

	return []ErrorStruct{
		errorFromErr(err),
	}
}

func errorFromErr(err error) ErrorStruct {
	if err == nil {
		msg := "error nil"
		return ErrorStruct{
			Message: &msg,
		}
	}

	msg := err.Error()
	return ErrorStruct{
		Message: &msg,
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "validateEmail":
		return "Should be email address"
	case "min":
		return "Should be at least " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	case "validateUsername":
		return "Should be greater than 3 and less than 50"
	}
	return "Unknown error"
}

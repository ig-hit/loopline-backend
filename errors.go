package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Error struct {
	ResponseCode int            `json:"-"`
	Description  string         `json:"description"`
	Details      []ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message     string `json:"message"`
	MessageCode string `json:"message_code"`
}

var (
	ErrInvalidForm = echo.NewHTTPError(http.StatusBadRequest, "invalid form")
	ErrInvalidToken       = echo.NewHTTPError(http.StatusForbidden, "invalid token")
	ErrNotFound              = echo.NewHTTPError(http.StatusNotFound, "not found")
)

func (e Error) Error() string {
	return e.Description
}

var (
	InternalServerError = &Error{
		ResponseCode: http.StatusInternalServerError,
		Description:  http.StatusText(http.StatusInternalServerError),
		Details: []ErrorDetails{
			{
				Message:     http.StatusText(http.StatusInternalServerError),
				MessageCode: "internal_server_error",
			},
		},
	}

	NotFoundError = &Error{
		ResponseCode: http.StatusNotFound,
		Description:  "Resource not found",
		Details: []ErrorDetails{
			{
				Message:     "Resource not found",
				MessageCode: "not_found",
			},
		},
	}
)

func convertHttpError(e echo.HTTPError) *Error {
	message := http.StatusText(http.StatusInternalServerError)

	switch e.Message.(type) {
	case string:
		message = strings.ToLower(e.Message.(string))
	}

	return &Error{
		ResponseCode: e.Code,
		Description:  message,
		Details: []ErrorDetails{
			{
				Message:     message,
				MessageCode: message,
			},
		},
	}
}

func convertValidationError(e validation.Errors) *Error {
	var details []ErrorDetails
	for field, fieldError := range e {
		details = append(details, ErrorDetails{
			Message:     fieldError.Error(),
			MessageCode: "invalid." + field,
		})
	}

	return &Error{
		ResponseCode: http.StatusBadRequest,
		Description:  "Validation error",
		Details:      details,
	}
}

func ErrorHandler(err error, c echo.Context) {
	var exception *Error

	switch err := err.(type) {
	case *echo.HTTPError:
		exception = convertHttpError(*err)
	case validation.Errors:
		exception = convertValidationError(err)
	default:
		switch err.Error() {
		case "not found":
			exception = NotFoundError
		default:
			exception = InternalServerError
			exception.Details = []ErrorDetails{
				{
					Message: err.Error(),
				},
			}
		}
	}

	resCode := c.Response().Status
	if exception.ResponseCode == 500 && resCode >= 400 {
		exception.ResponseCode = resCode
	}

	_ = c.JSON(exception.ResponseCode, exception)
}

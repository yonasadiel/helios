package helios

import "net/http"

// Error is interface of error that can be thrown
// for response
type Error interface {
	GetMessage() map[string]interface{}
	GetStatusCode() int
}

// APIError is standardized error of Charon app
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

// GetMessage returns the message to shown as response body
func (apiError APIError) GetMessage() map[string]interface{} {
	message := make(map[string]interface{})
	message["code"] = apiError.Code
	message["message"] = apiError.Message

	return message
}

// GetStatusCode returns the http status code
func (apiError APIError) GetStatusCode() int {
	return apiError.StatusCode
}

// FormError is common error, usually after parsing the request body
type FormError struct {
	Code          string
	FieldError    map[string]([]string)
	NonFieldError []string
}

// AddFieldError pushes the errorMessage to the field's error list
func (formError *FormError) AddFieldError(fieldName string, errorMessage string) {
	if len(formError.FieldError) == 0 {
		formError.FieldError = make(map[string]([]string))
	}
	if _, ok := formError.FieldError[fieldName]; !ok {
		formError.FieldError[fieldName] = make([]string, 0)
	}
	formError.FieldError[fieldName] = append(formError.FieldError[fieldName], errorMessage)
}

// AddNonFieldError pushes the errorMessage to the nonfield error list
func (formError *FormError) AddNonFieldError(errorMessage string) {
	formError.NonFieldError = append(formError.NonFieldError, errorMessage)
}

// GetFieldErrors returns the copy of message to shown as response body.
// It is dictionary of string to list of errors, with field name as the key
func (formError FormError) GetFieldErrors() map[string]([]string) {
	fieldsError := make(map[string]([]string))
	for k, v := range formError.FieldError {
		fieldError := make([]string, len(v))
		copy(fieldError, v)
		fieldsError[k] = fieldError
	}
	return fieldsError
}

// GetNonFieldErrors returns the copy of non field errors
func (formError FormError) GetNonFieldErrors() []string {
	nonFieldError := make([]string, len(formError.NonFieldError))
	copy(nonFieldError, formError.NonFieldError)
	return nonFieldError
}

// GetMessage returns the message to shown as response body
// it will include code (unique identifier) and map as message
// the message will contain field name as key and error as value
func (formError FormError) GetMessage() map[string]interface{} {
	messageFields := formError.GetFieldErrors()
	messageFields["_error"] = formError.GetNonFieldErrors()

	message := make(map[string]interface{})
	if formError.Code == "" {
		message["code"] = "form_error"
	} else {
		message["code"] = formError.Code
	}
	message["message"] = messageFields
	return message
}

// IsError returns true if there is at least one error,
// and return false if there is no error on the struct
func (formError FormError) IsError() bool {
	return len(formError.FieldError) > 0 || len(formError.NonFieldError) > 0
}

// GetStatusCode returns HTTP 400 Bad Request code
func (formError FormError) GetStatusCode() int {
	return http.StatusBadRequest
}

// ErrInternalServerError is general error that will be send
// if there is unexpected error on the server
var ErrInternalServerError = APIError{
	StatusCode: http.StatusInternalServerError,
	Code:       "internal_server_error",
	Message:    "Error occured while processing the request",
}

// ErrUnsupportedContentType returned when request has content-type
// header that is unsupported
var ErrUnsupportedContentType = APIError{
	StatusCode: http.StatusUnsupportedMediaType,
	Code:       "unsupported_content_type",
	Message:    "Currently, we are accepting application/json only",
}

// ErrJSONParseFailed will be returned when calling req.DeserializeRequestData
// but with bad JSON. For example, int field that supplied with string
var ErrJSONParseFailed = APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "failed_to_parse_json",
	Message:    "Failed to parse json request",
}

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
	FieldError    map[string]([]string)
	NonFieldError []string
}

// AddFieldError pushes the errorMessage to the field's error list
func (formError FormError) AddFieldError(fieldName string, errorMessage string) {
	if _, ok := formError.FieldError[fieldName]; !ok {
		formError.FieldError[fieldName] = make([]string, 0)
	}
	formError.FieldError[fieldName] = append(formError.FieldError[fieldName], errorMessage)
}

// GetMessage returns the message to shown as response body
// it will include code (unique identifier) and map as message
// the message will contain field name as key and error as value
func (formError FormError) GetMessage() map[string]interface{} {
	messageFields := make(map[string]([]string))
	for k, v := range formError.FieldError {
		fieldErrors := make([]string, len(v))
		copy(fieldErrors, v)
		messageFields[k] = fieldErrors
	}
	nonFieldError := make([]string, len(formError.NonFieldError))
	copy(nonFieldError, formError.NonFieldError)
	messageFields["_error"] = nonFieldError

	message := make(map[string]interface{})
	message["code"] = "form_error"
	message["message"] = messageFields
	return message
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

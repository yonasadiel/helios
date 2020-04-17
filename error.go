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
	FieldError    NestedFieldError
	NonFieldError AtomicFieldError
}

// FieldError is the interface of error for a field in a form. A field
// in a form can be:
// - an atomic field (ex: username, fullname), using AtomicFieldError
// - an array field (ex: todos, checkboxes), using ArrayFieldError
// - a nested field (ex: address consist of zip code, city), using NestedFieldError
type FieldError interface {
	GetMessage() interface{}
	IsError() bool
}

// AtomicFieldError is error representation of one field, example:
// passwordErr := AtomicFieldError{"password is too short". "password must include symbol"}
// will converted to json:
// ["password is too short", "password must include symbol"]
type AtomicFieldError []string

// GetMessage returns the json-friendly array copy of the error
func (err AtomicFieldError) GetMessage() interface{} {
	message := make([]string, 0)
	for _, e := range err {
		message = append(message, e)
	}
	return message
}

// IsError returns true if there is any error
func (err AtomicFieldError) IsError() bool {
	return len(err) > 0
}

// ArrayFieldError is error representation of array field, example:
// todosErr := ArrayFieldError{AtomicFieldError{"todo can't be empty"}, AtomicFieldError{}, AtomicFieldError{"todo is duplicate"}}
// will converted to json:
// [["todo can't be empty"],[],["todo is duplicate"]]
type ArrayFieldError []FieldError

// GetMessage returns json-friendly array copy of the error
func (err ArrayFieldError) GetMessage() interface{} {
	message := make([]interface{}, 0)
	for _, e := range err {
		message = append(message, e.GetMessage())
	}
	return message
}

// IsError returns true if there is at least one member with err.
// Tt will iterate all the member.
func (err ArrayFieldError) IsError() bool {
	var isError bool = false
	for _, e := range err {
		if e.IsError() {
			isError = true
		}
	}
	return isError
}

// NestedFieldError is error representation other field error mapped by field name, example:
// addressErr := NestedFieldError{
//   "zipCode": AtomicFieldError{"zip code is invalid"},
//   "city": AtomicFieldError{"city can't be empty"},
// }
// will converted to json:
// {"city":["city can't be empty"],"zipCode":["zip code is invalid"]}
type NestedFieldError map[string]FieldError

// GetMessage returns json-friendly map copy of the error
func (err NestedFieldError) GetMessage() interface{} {
	message := make(map[string]interface{})
	for k, v := range err {
		message[k] = v.GetMessage()
	}
	return message
}

// IsError returns true if there is at least one member with err.
// It will iterate all the member.
func (err NestedFieldError) IsError() bool {
	var isError bool = false
	for _, v := range err {
		if v.IsError() {
			isError = true
		}
	}
	return isError
}

// GetMessage returns the message to shown as response body
// it will include code (unique identifier) and map as message
// the message will contain field name as key and error as value
func (formError FormError) GetMessage() map[string]interface{} {
	messageFields := make(map[string]interface{})
	for k, v := range formError.FieldError {
		messageFields[k] = v.GetMessage()
	}
	messageFields["_error"] = formError.NonFieldError.GetMessage()

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
	return formError.FieldError.IsError() || formError.NonFieldError.IsError()
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

package helios

import "net/http"

// Error is interface of error that can be thrown
// for response
type Error interface {
	GetMessage() map[string]interface{}
	GetStatusCode() int
}

// ErrorAPI is standardized error of Charon app
type ErrorAPI struct {
	StatusCode int
	Code       string
	Message    string
}

// GetMessage returns the message to shown as response body
func (apiError ErrorAPI) GetMessage() map[string]interface{} {
	message := make(map[string]interface{})
	message["code"] = apiError.Code
	message["message"] = apiError.Message

	return message
}

// GetStatusCode returns the http status code
func (apiError ErrorAPI) GetStatusCode() int {
	return apiError.StatusCode
}

// ErrorForm is common error, usually after parsing the request body
type ErrorForm struct {
	Code              string
	ErrorFormField    ErrorFormFieldNested
	NonErrorFormField ErrorFormFieldAtomic
}

// ErrorFormField is the interface of error for a field in a form. A field
// in a form can be:
// - an atomic field (ex: username, fullname), using ErrorFormFieldAtomic
// - an array field (ex: todos, checkboxes), using ErrorFormFieldArray
// - a nested field (ex: address consist of zip code, city), using ErrorFormFieldNested
type ErrorFormField interface {
	GetMessage() interface{}
	IsError() bool
}

// ErrorFormFieldAtomic is error representation of one field, example:
//     passwordErr := ErrorFormFieldAtomic{"password is too short". "password must include symbol"}
//     // will be converted to json:
//     ["password is too short","password must include symbol"]
type ErrorFormFieldAtomic []string

// GetMessage returns the json-friendly array copy of the error
func (err ErrorFormFieldAtomic) GetMessage() interface{} {
	message := make([]string, 0)
	for _, e := range err {
		message = append(message, e)
	}
	return message
}

// IsError returns true if there is any error
func (err ErrorFormFieldAtomic) IsError() bool {
	return len(err) > 0
}

// ErrorFormFieldArray is error representation of array field, example:
//     todosErr := ErrorFormFieldArray{ErrorFormFieldAtomic{"todo can't be empty"}, ErrorFormFieldAtomic{}, ErrorFormFieldAtomic{"todo is duplicate"}}
//     // will be converted to json:
//     [["todo can't be empty"],[],["todo is duplicate"]]
type ErrorFormFieldArray []ErrorFormField

// GetMessage returns json-friendly array copy of the error
func (err ErrorFormFieldArray) GetMessage() interface{} {
	message := make([]interface{}, 0)
	for _, e := range err {
		message = append(message, e.GetMessage())
	}
	return message
}

// IsError returns true if there is at least one member with err.
// Tt will iterate all the member.
func (err ErrorFormFieldArray) IsError() bool {
	var isError bool = false
	for _, e := range err {
		if e.IsError() {
			isError = true
		}
	}
	return isError
}

// ErrorFormFieldNested is error representation other field error mapped by field name, example:
//     addressErr := ErrorFormFieldNested{
//         "zipCode": ErrorFormFieldAtomic{"zip code is invalid"},
//         "city": ErrorFormFieldAtomic{"city can't be empty"},
//     }
//     // will be converted to json:
//     {"city":["city can't be empty"],"zipCode":["zip code is invalid"]}
type ErrorFormFieldNested map[string]ErrorFormField

// GetMessage returns json-friendly map copy of the error
func (err ErrorFormFieldNested) GetMessage() interface{} {
	message := make(map[string]interface{})
	for k, v := range err {
		message[k] = v.GetMessage()
	}
	return message
}

// IsError returns true if there is at least one member with err.
// It will iterate all the member.
func (err ErrorFormFieldNested) IsError() bool {
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
func (formError ErrorForm) GetMessage() map[string]interface{} {
	messageFields := make(map[string]interface{})
	for k, v := range formError.ErrorFormField {
		messageFields[k] = v.GetMessage()
	}
	messageFields["_error"] = formError.NonErrorFormField.GetMessage()

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
func (formError ErrorForm) IsError() bool {
	return formError.ErrorFormField.IsError() || formError.NonErrorFormField.IsError()
}

// GetStatusCode returns HTTP 400 Bad Request code
func (formError ErrorForm) GetStatusCode() int {
	return http.StatusBadRequest
}

// ErrInternalServerError is general error that will be send
// if there is unexpected error on the server
var ErrInternalServerError = ErrorAPI{
	StatusCode: http.StatusInternalServerError,
	Code:       "internal_server_error",
	Message:    "Error occured while processing the request",
}

// ErrUnsupportedContentType returned when request has content-type
// header that is unsupported
var ErrUnsupportedContentType = ErrorAPI{
	StatusCode: http.StatusUnsupportedMediaType,
	Code:       "unsupported_content_type",
	Message:    "Currently, we are accepting application/json only",
}

// ErrJSONParseFailed will be returned when calling req.DeserializeRequestData
// but with bad JSON. For example, int field that supplied with string
var ErrJSONParseFailed = ErrorAPI{
	StatusCode: http.StatusBadRequest,
	Code:       "failed_to_parse_json",
	Message:    "Failed to parse json request",
}

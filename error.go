package helios

import "net/http"

// APIError is standardized error of Charon app
type APIError struct {
	error
	StatusCode int
	Code       string
	Message    string
}

// GetMessage Get the message to shown as response body
func (apiError APIError) GetMessage() map[string]string {
	message := make(map[string]string)
	message["code"] = apiError.Code
	message["message"] = apiError.Message

	return message
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

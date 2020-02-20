package helios

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Request interface of Helios Http Request Wrapper
type Request interface {
	GetRequestData() map[string]string
	GetContextData(key string) interface{}
	GetSessionData(key string) interface{}

	SetContextData(key string, value interface{})
	SetSessionData(key string, value interface{})
	SaveSession()

	SendJSON(output interface{}, code int)
}

// HTTPHandler receive Helios wrapped request and ressponse
type HTTPHandler func(Request)

// HTTPRequest wrapper of Helios Http Request
// r is the HTTP Request, containing request data
// w is the HTTP Response writer, to write HTTP reply
// s is the session of current request, using gorilla/sessions package
// c is the context of the current request, can be used for user data, etc
type HTTPRequest struct {
	r *http.Request
	w http.ResponseWriter
	s *sessions.Session
	c map[string]interface{}
}

// GetRequestData return the data of request
func (req *HTTPRequest) GetRequestData() map[string]string {
	return mux.Vars(req.r)
}

// GetSessionData return the data of session with known key
func (req *HTTPRequest) GetSessionData(key string) interface{} {
	return req.s.Values[key]
}

// SetSessionData set the data of session
func (req *HTTPRequest) SetSessionData(key string, value interface{}) {
	req.s.Values[key] = value
}

// GetContextData return the data of session with known key
func (req *HTTPRequest) GetContextData(key string) interface{} {
	return req.c[key]
}

// SetContextData set the data of session
func (req *HTTPRequest) SetContextData(key string, value interface{}) {
	req.c[key] = value
}

// SaveSession saves the session to the cookie
func (req *HTTPRequest) SaveSession() {
	req.s.Save(req.r, req.w)
}

// SendJSON write json as http response
func (req *HTTPRequest) SendJSON(output interface{}, code int) {
	response, _ := json.Marshal(output)

	req.w.Header().Set("Content-Type", "application/json")
	req.w.WriteHeader(code)
	req.w.Write(response)
}

// MockRequest is Request object that is mocked for testing purposes
type MockRequest struct {
	RequestData  map[string]string
	SessionData  map[string]interface{}
	ContextData  map[string]interface{}
	JSONResponse []byte
	StatusCode   int
}

// GetRequestData return the data of request
func (req *MockRequest) GetRequestData() map[string]string {
	return req.RequestData
}

// SetRequestData set the data of session
func (req *MockRequest) SetRequestData(key string, value string) {
	req.RequestData[key] = value
}

// GetSessionData return the data of session with known key
func (req *MockRequest) GetSessionData(key string) interface{} {
	return req.SessionData[key]
}

// SetSessionData set the data of session
func (req *MockRequest) SetSessionData(key string, value interface{}) {
	req.SessionData[key] = value
}

// GetContextData return the data of session with known key
func (req *MockRequest) GetContextData(key string) interface{} {
	return req.ContextData[key]
}

// SetContextData set the data of session
func (req *MockRequest) SetContextData(key string, value interface{}) {
	req.ContextData[key] = value
}

// SaveSession do nothing because the session is already saved
func (req *MockRequest) SaveSession() {
	//
}

// SendJSON write json as http response
func (req *MockRequest) SendJSON(output interface{}, code int) {
	var err error
	req.JSONResponse, err = json.Marshal(output)
	if err != nil {
		req.StatusCode = http.StatusInternalServerError
	} else {
		req.StatusCode = code
	}
}

// NewMockRequest returns new MockRequest with empty data
func NewMockRequest() MockRequest {
	return MockRequest{
		SessionData: make(map[string]interface{}),
		RequestData: make(map[string]string),
		ContextData: make(map[string]interface{}),
	}
}

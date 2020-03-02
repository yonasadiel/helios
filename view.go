package helios

import (
	"encoding/json"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Request interface of Helios Http Request Wrapper
type Request interface {
	DeserializeRequestData(obj interface{}) *APIError

	GetURLParam(key string) string

	GetContextData(key string) interface{}
	SetContextData(key string, value interface{})

	GetSessionData(key string) interface{}
	SetSessionData(key string, value interface{})
	SaveSession()

	ClientIP() string

	SendJSON(output interface{}, code int)
}

// HTTPHandler receive Helios wrapped request and ressponse
type HTTPHandler func(Request)

// HTTPRequest wrapper of Helios Http Request
// r is the HTTP Request, containing request data
// w is the HTTP Response writer, to write HTTP reply
// s is the session of current request, using gorilla/sessions package
// c is the context of the current request, can be used for user data, etc
// u is the url params argument
type HTTPRequest struct {
	r *http.Request
	w http.ResponseWriter
	s *sessions.Session
	c map[string]interface{}
	u map[string]string
}

// NewHTTPRequest wraps usual http request and response writer to HTTPRequest struct
func NewHTTPRequest(w http.ResponseWriter, r *http.Request) HTTPRequest {
	return HTTPRequest{
		r: r,
		w: w,
		s: App.getSession(r),
		c: make(map[string]interface{}),
		u: mux.Vars(r),
	}
}

// GetURLParam return the parameter of the request url
func (req *HTTPRequest) GetURLParam(key string) string {
	return req.u[key]
}

// DeserializeRequestData deserializes the request body
// and parse it into pointer to struct
func (req *HTTPRequest) DeserializeRequestData(obj interface{}) *APIError {
	contentType := req.r.Header.Get("Content-Type")
	if contentType == "application/json" || contentType == "" {
		decoder := json.NewDecoder(req.r.Body)
		err := decoder.Decode(obj)
		if err != nil {
			return &ErrUnsupportedContentType
		}
		return nil
	}
	return &ErrUnsupportedContentType
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
	req.s.Save(req.r, req.w) // nolint:errcheck
}

// SendJSON write json as http response
func (req *HTTPRequest) SendJSON(output interface{}, code int) {
	response, _ := json.Marshal(output)

	req.w.Header().Set("Content-Type", "application/json")
	req.w.WriteHeader(code)
	req.w.Write(response) // nolint:errcheck
}

// ClientIP returns the original ip address of the request.
// First, it checks for X-Forwarded-For and X-Real-Ip http header
// (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For)
// If they are not present, return the http.Request.RemoteAddr
// The priority is: X-Forwarded-For, X-Real-Ip, RemoteAddr
func (req *HTTPRequest) ClientIP() string {
	clientIP := req.r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(req.r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// MockRequest is Request object that is mocked for testing purposes
type MockRequest struct {
	RequestData  interface{}
	SessionData  map[string]interface{}
	ContextData  map[string]interface{}
	JSONResponse []byte
	StatusCode   int
	URLParam     map[string]string
	RemoteAddr   string
}

// NewMockRequest returns new MockRequest with empty data
// RemoteAddr is set to 127.0.0.1 in default
func NewMockRequest() MockRequest {
	return MockRequest{
		SessionData: make(map[string]interface{}),
		RequestData: make(map[string]string),
		ContextData: make(map[string]interface{}),
		URLParam:    make(map[string]string),
		RemoteAddr:  "127.0.0.1",
	}
}

// GetURLParam returns the url param of given key
func (req *MockRequest) GetURLParam(key string) string {
	return req.URLParam[key]
}

// DeserializeRequestData return the data of request
func (req *MockRequest) DeserializeRequestData(obj interface{}) *APIError {
	if req.RequestData == nil {
		return &ErrUnsupportedContentType
	}
	result := reflect.ValueOf(obj).Elem()
	result.Set(reflect.ValueOf(req.RequestData))
	return nil
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

// ClientIP returns RemoteAddr data of req
func (req *MockRequest) ClientIP() string {
	return req.RemoteAddr
}

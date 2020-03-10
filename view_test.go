package helios

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleRequest struct {
	A string `json:"a"`
	B int    `json:"b"`
	C bool   `json:"c"`
	D string `json:"d"`
	E int    `json:"e"`
	F bool   `json:"f"`
}

func TestHandle(t *testing.T) {
	App.BeforeTest()

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/def", nil)

	f := func(req Request) {
		json := make(map[string]int)
		json["abc"] = 2
		json["def"] = 3
		req.SendJSON(json, 201)
	}

	Handle(f)(recorder, request)
	actualResponse := make([]byte, 17)
	expectedResponse := []byte("{\"abc\":2,\"def\":3}")
	n, err := recorder.Result().Body.Read(actualResponse)

	assert.Nil(t, err, "Fail on reading body")
	assert.Equal(t, n, len(expectedResponse), "Different number of bytes")
	assert.Equal(t, expectedResponse, actualResponse, "Different body")
}

func TestMockRequest(t *testing.T) {
	App.BeforeTest()

	req := NewMockRequest()

	req.URLParam["abc"] = "def"
	req.URLParam["one"] = "1"
	req.URLParam["long"] = "123456789012345678901234567890"
	assert.Equal(t, "def", req.GetURLParam("abc"), "Failed to retrieve url param")
	assert.Equal(t, "", req.GetURLParam("def"), "Not found url param should return empty string")
	oneParam, errOneParam := req.GetURLParamUint("one")
	_, errTwoParam := req.GetURLParamUint("two")
	_, errAbcParam := req.GetURLParamUint("abc")
	_, errLongParam := req.GetURLParamUint("long")
	assert.Equal(t, uint(1), oneParam, "Different result on uint url param")
	assert.Nil(t, errOneParam, "Failed to parse uint url param")
	assert.NotNil(t, errTwoParam, "Empty uint param will return error")
	assert.NotNil(t, errAbcParam, "Not number param will return error")
	assert.NotNil(t, errLongParam, "Long number param will return error")

	var expected sampleRequest = sampleRequest{A: "def"}
	var actual sampleRequest

	req.RequestData = nil
	errDeserialize1 := req.DeserializeRequestData(&actual)
	assert.NotNil(t, errDeserialize1, "If request data is nil, DeserializeRequestData purposely throw error")

	req.RequestData = expected
	errDeserialize2 := req.DeserializeRequestData(&actual)
	assert.Nil(t, errDeserialize2, "Failed to deserialize request data")
	assert.Equal(t, expected, actual, "Differenet request data")

	req.SetContextData("abc", 2)
	req.SetContextData("abc", 3)
	req.SetContextData("def", "ghi")

	assert.Equal(t, 3, req.GetContextData("abc"), "Fail to use context data")
	assert.Equal(t, "ghi", req.GetContextData("def"), "Fail to use context data")
	assert.Nil(t, req.GetContextData("ghi"), "Context data with unexist key should return nil")

	req.SetSessionData("abc", 4)
	req.SetSessionData("abc", 5)
	req.SetSessionData("def", true)
	req.SaveSession()

	assert.Equal(t, 5, req.GetSessionData("abc"), "Fail to use session data")
	assert.Equal(t, true, req.GetSessionData("def"), "Fail to use session data")
	assert.Nil(t, req.GetSessionData("ghi"), "Session data with unexist key should return nil")

	req.RequestHeader["header-a"] = "a"
	assert.Empty(t, req.GetHeader("header-b"), "Missing header should return empty string")
	assert.Equal(t, "a", req.GetHeader("header-a"), "Different request header value")
	assert.Equal(t, "a", req.GetHeader("HEADER-A"), "Request header should be case insensitive")

	req.SetHeader("header-x", "x")
	req.SetHeader("header-x", "y")
	assert.Empty(t, req.ResponseHeader["header-y"], "Unset header should return empty string")
	assert.Equal(t, "y", req.ResponseHeader["header-x"], "Different header set")

	req.SendJSON(sampleRequest{A: "abcde", B: 2, C: true}, 499)
	expectedResponse := "{\"a\":\"abcde\",\"b\":2,\"c\":true,\"d\":\"\",\"e\":0,\"f\":false}"
	assert.Equal(t, expectedResponse, string(req.JSONResponse), "Different JSON Response")
	assert.Equal(t, 499, req.StatusCode, "Different Response status code")

	assert.Equal(t, "127.0.0.1", req.ClientIP(), "Default for ClientIP is 127.0.0.1")
	req.RemoteAddr = "1.2.3.4"
	assert.Equal(t, "1.2.3.4", req.ClientIP(), "ClientIP() should returns the RemoteAddr attribute")
}

func TestNewHTTPRequest(t *testing.T) {
	App.BeforeTest()

	recorder := httptest.NewRecorder()
	response := []byte("abc")
	request, _ := http.NewRequest("GET", "/def", nil)
	req := NewHTTPRequest(recorder, request)
	_, err := req.w.Write(response)
	assert.Nil(t, err, "Fail on writing response")
	actualResponse := make([]byte, 5)
	n, err := recorder.Result().Body.Read(actualResponse)
	assert.Nil(t, err, "Fail on reading body")
	assert.Equal(t, n, 3, "Different number of bytes")
	assert.Equal(t, []byte{0x61, 0x62, 0x63, 0x0, 0x0}, actualResponse, "Different body")
}

func TestHTTPRequest(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader("{\"a\":\"abcde\",\"b\":2,\"c\":true}"))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("header-a", "a")
	recorder := httptest.NewRecorder()

	urlParam := make(map[string]string)
	urlParam["id"] = "3"
	urlParam["abc"] = "def"
	urlParam["long"] = "123456789012345678901234567890"

	contextData := make(map[string]interface{})
	contextData["abc"] = "random_context_data"
	contextData["def"] = "def"

	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: contextData,
		u: urlParam,
	}
	req.SetContextData("def", "ghi")

	assert.Empty(t, req.GetHeader("header-b"), "Missing header should return empty string")
	assert.Equal(t, "a", req.GetHeader("header-a"), "Different request header returned")
	assert.Equal(t, "a", req.GetHeader("HEADER-A"), "Header should be case insensitive")

	req.SetHeader("header-x", "x")
	req.SetHeader("header-x", "y")
	assert.Empty(t, req.w.Header().Get("header-y"), "Empty header shoul be empty string")
	assert.Equal(t, "y", req.w.Header().Get("header-x"), "Different header set")

	assert.Equal(t, "3", req.GetURLParam("id"), "Failed to retrive url param")
	idParam, errIDParam := req.GetURLParamUint("id")
	_, errOtherIDParam := req.GetURLParamUint("id2")
	_, errAbcParam := req.GetURLParamUint("abc")
	_, errLongParam := req.GetURLParamUint("long")
	assert.Equal(t, uint(3), idParam, "Different result on uint url param")
	assert.Nil(t, errIDParam, "Failed to parse uint url param")
	assert.NotNil(t, errOtherIDParam, "Empty uint param will return error")
	assert.NotNil(t, errAbcParam, "Not number param will return error")
	assert.NotNil(t, errLongParam, "Long number param will return error")
	assert.Equal(t, "random_context_data", req.GetContextData("abc"), "Failed to retrive context data")
	assert.Equal(t, "ghi", req.GetContextData("def"), "Failed to retrive modified context data")
	assert.Nil(t, req.GetContextData("ghi"), "Missing context data should treated as nil")
}

func TestHTTPRequestJSONEncoded(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader("{\"a\":\"abcde\",\"b\":2,\"c\":true}"))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	var requestData sampleRequest

	err := req.DeserializeRequestData(&requestData)
	assert.Nil(t, err, "Fail on deserializing request data")
	assert.Equal(t, "abcde", requestData.A, "Different string request data")
	assert.Equal(t, 2, requestData.B, "Different integer request data")
	assert.Equal(t, true, requestData.C, "Different boolean request data")
	assert.Equal(t, "", requestData.D, "Different empty string request data")
	assert.Equal(t, 0, requestData.E, "Different empty integer request data")
	assert.Equal(t, false, requestData.F, "Different empty boolean request data")
}

func TestHTTPRequestJSONPoorlyEncoded(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader("{\"a\":\"abcde\",\"b\":2,\"c\":true"))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	var requestData sampleRequest

	err := req.DeserializeRequestData(&requestData)
	assert.Equal(t, &ErrUnsupportedContentType, err, "Poorly encoded JSON should return Unsupported Content Type error")
}

func TestHTTPRequestUrlFormEncoded(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader("a=abcde&b=2&c=true}"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	var requestData sampleRequest

	err := req.DeserializeRequestData(&requestData)
	assert.Equal(t, &ErrUnsupportedContentType, err, "application/x-www-form-urlencoded is not supported yet")
}

func TestHTTPRequestMultipartFormData(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader(`
	-----------------------------9051914041544843365972754266
	Content-Disposition: form-data; name="a"

	abcde
	-----------------------------9051914041544843365972754266
	Content-Disposition: form-data; name="b"

	2
	-----------------------------9051914041544843365972754266
	Content-Disposition: form-data; name="c"

	true

	-----------------------------9051914041544843365972754266-
	`))
	request.Header.Set("Content-Type", "multipart/form-data")
	recorder := httptest.NewRecorder()
	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	var requestData sampleRequest

	err := req.DeserializeRequestData(&requestData)
	assert.Equal(t, &ErrUnsupportedContentType, err, "multipart/form-data is not supported yet")
}

func TestHTTPRequestClientIP(t *testing.T) {
	App.BeforeTest()

	requestXFF, _ := http.NewRequest("POST", "/def", nil)
	requestXFF.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	requestXFF.Header.Set("X-Real-Ip", "11.22.33.44")
	requestXFF.RemoteAddr = "55.66.77.88:12345"

	reqXFF := HTTPRequest{
		r: requestXFF,
		w: nil,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	assert.Equal(t, "1.2.3.4", reqXFF.ClientIP(), "If X-Forwarded-For header is present, ClientIP should return the first entry of the header")

	requestXRI, _ := http.NewRequest("POST", "/def", nil)
	requestXRI.Header.Set("X-Forwarded-For", "")
	requestXRI.Header.Set("X-Real-Ip", "11.22.33.44")
	requestXRI.RemoteAddr = "55.66.77.88:12345"

	reqXRI := HTTPRequest{
		r: requestXRI,
		w: nil,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	assert.Equal(t, "11.22.33.44", reqXRI.ClientIP(), "If X-Forwarded-For header is not present and X-Real-Ip is present, ClientIP should return the X-Real-Ip")

	requestRemoteAddr, _ := http.NewRequest("POST", "/def", nil)
	requestRemoteAddr.Header.Set("X-Forwarded-For", "")
	requestRemoteAddr.Header.Set("X-Real-Ip", "")
	requestRemoteAddr.RemoteAddr = "55.66.77.88:12345"

	reqRemoteAddr := HTTPRequest{
		r: requestRemoteAddr,
		w: nil,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	assert.Equal(t, "55.66.77.88", reqRemoteAddr.ClientIP(), "If X-Forwarded-For and X-Real-Ip headers are not present, ClientIP should return RemoteAddr")

	requestBadIP, _ := http.NewRequest("POST", "/def", nil)
	requestBadIP.Header.Set("X-Forwarded-For", "")
	requestBadIP.Header.Set("X-Real-Ip", "")
	requestBadIP.RemoteAddr = "55.66.77.88"

	reqBadIP := HTTPRequest{
		r: requestBadIP,
		w: nil,
		s: nil,
		c: make(map[string]interface{}),
		u: make(map[string]string),
	}

	assert.Empty(t, reqBadIP.ClientIP(), "Bad IP format on RemoteAddr will return empty ip")
}

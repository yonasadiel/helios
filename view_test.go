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

func TestMockRequest(t *testing.T) {
	App.BeforeTest()

	req := NewMockRequest()

	req.SetURLParam("abc", "def")
	assert.Equal(t, "def", req.GetURLParam("abc"), "Failed to retrieve url param")
	assert.Equal(t, "", req.GetURLParam("def"), "Not found url param should return empty string")

	var expected sampleRequest = sampleRequest{A: "def"}
	var actual sampleRequest

	req.SetRequestData(expected)

	err := req.DeserializeRequestData(&actual)
	assert.Nil(t, err, "Failed to deserialize request data")
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

	req.SendJSON(sampleRequest{A: "abcde", B: 2, C: true}, 499)
	expectedResponse := "{\"a\":\"abcde\",\"b\":2,\"c\":true,\"d\":\"\",\"e\":0,\"f\":false}"
	assert.Equal(t, expectedResponse, string(req.JSONResponse), "Different JSON Response")
	assert.Equal(t, 499, req.StatusCode, "Different Response status code")
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
	recorder := httptest.NewRecorder()

	urlParam := make(map[string]string)
	urlParam["id"] = "3"

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

	assert.Equal(t, "3", req.GetURLParam("id"), "Failed to retrive url param")
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

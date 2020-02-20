package helios

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockRequest(t *testing.T) {
	App.BeforeTest()

	req := NewMockRequest()

	req.SetRequestData("abc", "def")
	req.SetRequestData("abc", "ghi")

	data := req.GetRequestData()
	expected := make(map[string]string)
	expected["abc"] = "ghi"
	assert.Equal(t, expected, data, "Differenet request data")

	req.SetContextData("abc", 2)
	req.SetContextData("abc", 3)
	req.SetContextData("def", "ghi")

	assert.Equal(t, 3, req.GetContextData("abc"), "Fail to use context data")
	assert.Equal(t, "ghi", req.GetContextData("def"), "Fail to use context data")
	assert.Nil(t, req.GetContextData("ghi"), "Context data with unexist key should return nil")

	req.SetSessionData("abc", 4)
	req.SetSessionData("abc", 5)
	req.SetSessionData("def", true)

	assert.Equal(t, 5, req.GetSessionData("abc"), "Fail to use session data")
	assert.Equal(t, true, req.GetSessionData("def"), "Fail to use session data")
	assert.Nil(t, req.GetSessionData("ghi"), "Session data with unexist key should return nil")
}

func TestHTTPRequest(t *testing.T) {
	App.BeforeTest()

	request, _ := http.NewRequest("POST", "/def", strings.NewReader("abc=ghi&def=jkl"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	recorder := httptest.NewRecorder()
	req := HTTPRequest{
		r: request,
		w: recorder,
		s: nil,
		c: make(map[string]interface{}),
	}

	data := req.GetRequestData()
	expected := make(map[string]string)
	expected["abc"] = "ghi"
	assert.Equal(t, expected, data, "Differenet request data")
}

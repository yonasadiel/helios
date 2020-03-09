package helios

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareOrder(t *testing.T) {
	App.BeforeTest()

	order := make([]string, 3)
	nOrder := 0
	m1 := func(f HTTPHandler) HTTPHandler {
		return func(req Request) {
			order[nOrder] = "m1"
			nOrder++
			f(req)
		}
	}
	m2 := func(f HTTPHandler) HTTPHandler {
		return func(req Request) {
			order[nOrder] = "m2"
			nOrder++
			f(req)
		}
	}
	m3 := func(f HTTPHandler) HTTPHandler {
		return func(req Request) {
			order[nOrder] = "m3"
			nOrder++
			f(req)
		}
	}
	f := func(req Request) {
		// empty handler
	}
	fm := makeMiddleware(f, []Middleware{m1, m2, m3})
	req := NewMockRequest()
	fm(&req)

	assert.Equal(t, "m1", order[0], "First middleware to be executed should be m1")
	assert.Equal(t, "m2", order[1], "First middleware to be executed should be m2")
	assert.Equal(t, "m3", order[2], "First middleware to be executed should be m3")
}

func TestWithMiddleware(t *testing.T) {
	App.BeforeTest()

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/def", nil)

	f := func(req Request) {
		json := make(map[string]int)
		json["abc"] = 2
		json["def"] = 3
		req.SendJSON(json, 201)
	}

	WithMiddleware(f, []Middleware{})(recorder, request)
	actualResponse := make([]byte, 17)
	expectedResponse := []byte("{\"abc\":2,\"def\":3}")
	n, err := recorder.Result().Body.Read(actualResponse)

	assert.Nil(t, err, "Fail on reading body")
	assert.Equal(t, n, len(expectedResponse), "Different number of bytes")
	assert.Equal(t, expectedResponse, actualResponse, "Different body")
}

func TestCreateCORSMiddleware(t *testing.T) {
	strictCORS := CreateCORSMiddleware([]string{"http://localhost:9001", "http://localhost:9002"})
	wildcardCORS := CreateCORSMiddleware([]string{"*"})

	request, _ := http.NewRequest("GET", "/def", nil)
	f := func(req Request) {
		json := make(map[string]int)
		json["abc"] = 2
		json["def"] = 3
		req.SendJSON(json, 200)
	}

	recorder1 := httptest.NewRecorder()
	request.Header.Set("Origin", "http://localhost:9001")
	WithMiddleware(f, []Middleware{strictCORS})(recorder1, request)
	assert.Equal(t, "http://localhost:9001", recorder1.Header().Get("Access-Control-Allow-Origin"), "Fail to allow an origin")

	recorder2 := httptest.NewRecorder()
	request.Header.Set("Origin", "http://localhost:9002")
	WithMiddleware(f, []Middleware{strictCORS})(recorder2, request)
	assert.Equal(t, "http://localhost:9002", recorder2.Header().Get("Access-Control-Allow-Origin"), "Fail to allow another origin")

	recorder3 := httptest.NewRecorder()
	request.Header.Set("Origin", "http://localhost:9003")
	WithMiddleware(f, []Middleware{strictCORS})(recorder3, request)
	assert.Equal(t, "", recorder3.Header().Get("Access-Control-Allow-Origin"), "Fail to disallow an origin")

	recorder4 := httptest.NewRecorder()
	request.Header.Set("Origin", "http://localhost:9003")
	WithMiddleware(f, []Middleware{wildcardCORS})(recorder4, request)
	assert.Equal(t, "http://localhost:9003", recorder4.Header().Get("Access-Control-Allow-Origin"), "Fail to allow an origin on wildcard CORS")
}

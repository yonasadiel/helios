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

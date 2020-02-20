package helios

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	App.BeforeTest()

	err := APIError{
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
		Message:    "Not Found",
	}

	msg := err.GetMessage()

	assert.Equal(t, "not_found", msg["code"], "Wrong code")
	assert.Equal(t, "Not Found", msg["message"], "Wrong message")
}

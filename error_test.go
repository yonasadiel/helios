package helios

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	App.BeforeTest()

	var err Error = APIError{
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
		Message:    "Not Found",
	}

	msg := err.GetMessage()

	assert.Equal(t, http.StatusNotFound, err.GetStatusCode(), "Wrong status code")
	assert.Equal(t, "not_found", msg["code"], "Wrong code")
	assert.Equal(t, "Not Found", msg["message"], "Wrong message")
}

func TestFormError(t *testing.T) {
	App.BeforeTest()

	var err FormError = FormError{
		FieldError:    make(map[string]([]string)),
		NonFieldError: []string{"err1", "err2"},
	}
	err.AddFieldError("field1", "err3")
	err.AddFieldError("field1", "err4")
	err.AddFieldError("field2", "err5")
	var err2 Error = err

	assert.Equal(t, http.StatusBadRequest, err.GetStatusCode(), "Qrong status code")
	res, err3 := json.Marshal(err2.GetMessage())
	assert.Nil(t, err3, "Failed to json marshal")
	assert.Equal(t, `{"code":"form_error","message":{"_error":["err1","err2"],"field1":["err3","err4"],"field2":["err5"]}}`, string(res), "Different message")
}

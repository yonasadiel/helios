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

	var err FormError = FormError{}
	err.NonFieldError = []string{"err1", "err2"}
	err.AddFieldError("field1", "err3")
	err.AddFieldError("field1", "err4")
	err.AddFieldError("field2", "err5")
	var errCasted Error
	errCasted = err

	assert.Equal(t, http.StatusBadRequest, err.GetStatusCode(), "Qrong status code")
	var marshalError error
	var res []byte

	res, marshalError = json.Marshal(err.GetFieldErrors())
	assert.Nil(t, marshalError, "Failed to json marshal")
	assert.Equal(t, `{"field1":["err3","err4"],"field2":["err5"]}`, string(res), "Different message")

	res, marshalError = json.Marshal(err.GetNonFieldErrors())
	assert.Nil(t, marshalError, "Failed to json marshal")
	assert.Equal(t, `["err1","err2"]`, string(res), "Different message")

	res, marshalError = json.Marshal(errCasted.GetMessage())
	assert.Nil(t, marshalError, "Failed to json marshal")
	assert.Equal(t, `{"code":"form_error","message":{"_error":["err1","err2"],"field1":["err3","err4"],"field2":["err5"]}}`, string(res), "Different message")

	err.Code = "changed_error_code"
	errCasted = err
	res, marshalError = json.Marshal(errCasted.GetMessage())
	assert.Nil(t, marshalError, "Failed to json marshal")
	assert.Equal(t, `{"code":"changed_error_code","message":{"_error":["err1","err2"],"field1":["err3","err4"],"field2":["err5"]}}`, string(res), "Different message")
}

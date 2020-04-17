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

	type formErrorTestCase struct {
		err                FormError
		expectedStatusCode int
		expectedJSON       string
		expectedIsError    bool
	}
	testCases := []formErrorTestCase{{
		err:                FormError{},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[]}}`,
		expectedIsError:    false,
	}, {
		err: FormError{
			FieldError: NestedFieldError{
				"field1": AtomicFieldError{},
				"field2": ArrayFieldError{
					AtomicFieldError{},
					AtomicFieldError{},
				},
			},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[],"field1":[],"field2":[[],[]]}}`,
		expectedIsError:    false,
	}, {
		err:                FormError{Code: "custom_code"},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"custom_code","message":{"_error":[]}}`,
		expectedIsError:    false,
	}, {
		err: FormError{
			NonFieldError: []string{"err1", "err2"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":["err1","err2"]}}`,
		expectedIsError:    true,
	}, {
		err: FormError{
			FieldError: NestedFieldError{
				"atomic": AtomicFieldError{"err1", "err2"},
				"array": ArrayFieldError{
					NestedFieldError{
						"field1": AtomicFieldError{"err3"},
						"field2": ArrayFieldError{
							AtomicFieldError{"err4", "err5"},
						},
					},
					NestedFieldError{},
					NestedFieldError{
						"field1": AtomicFieldError{"err6"},
					},
				},
			},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[],"array":[{"field1":["err3"],"field2":[["err4","err5"]]},{},{"field1":["err6"]}],"atomic":["err1","err2"]}}`,
		expectedIsError:    true,
	}}
	for i, testCase := range testCases {
		t.Logf("TestFormError testcase #%d", i)
		var jsonRepresentation []byte
		var errMashalling error
		var err Error = testCase.err // cast the FormError to Error
		jsonRepresentation, errMashalling = json.Marshal(err.GetMessage())
		assert.Nil(t, errMashalling)
		assert.Equal(t, testCase.expectedStatusCode, err.GetStatusCode())
		assert.Equal(t, testCase.expectedIsError, testCase.err.IsError())
		assert.Equal(t, testCase.expectedJSON, string(jsonRepresentation))
	}
}

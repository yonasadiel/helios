package helios

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorAPI(t *testing.T) {
	App.BeforeTest()

	var err Error = ErrorAPI{
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
		Message:    "Not Found",
	}

	msg := err.GetMessage()

	assert.Equal(t, http.StatusNotFound, err.GetStatusCode(), "Wrong status code")
	assert.Equal(t, "not_found", msg["code"], "Wrong code")
	assert.Equal(t, "Not Found", msg["message"], "Wrong message")
}

func TestErrorForm(t *testing.T) {
	App.BeforeTest()

	type formErrorTestCase struct {
		err                ErrorForm
		expectedStatusCode int
		expectedJSON       string
		expectedIsError    bool
	}
	testCases := []formErrorTestCase{{
		err:                ErrorForm{},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[]}}`,
		expectedIsError:    false,
	}, {
		err: ErrorForm{
			ErrorFormField: ErrorFormFieldNested{
				"field1": ErrorFormFieldAtomic{},
				"field2": ErrorFormFieldArray{
					ErrorFormFieldAtomic{},
					ErrorFormFieldAtomic{},
				},
			},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[],"field1":[],"field2":[[],[]]}}`,
		expectedIsError:    false,
	}, {
		err:                ErrorForm{Code: "custom_code"},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"custom_code","message":{"_error":[]}}`,
		expectedIsError:    false,
	}, {
		err: ErrorForm{
			NonErrorFormField: []string{"err1", "err2"},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":["err1","err2"]}}`,
		expectedIsError:    true,
	}, {
		err: ErrorForm{
			ErrorFormField: ErrorFormFieldNested{
				"atomic": ErrorFormFieldAtomic{"err1", "err2"},
				"array": ErrorFormFieldArray{
					ErrorFormFieldNested{
						"field1": ErrorFormFieldAtomic{"err3"},
						"field2": ErrorFormFieldArray{
							ErrorFormFieldAtomic{"err4", "err5"},
						},
					},
					ErrorFormFieldNested{},
					ErrorFormFieldNested{
						"field1": ErrorFormFieldAtomic{"err6"},
					},
				},
			},
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedJSON:       `{"code":"form_error","message":{"_error":[],"array":[{"field1":["err3"],"field2":[["err4","err5"]]},{},{"field1":["err6"]}],"atomic":["err1","err2"]}}`,
		expectedIsError:    true,
	}}
	for i, testCase := range testCases {
		t.Logf("TestErrorForm testcase #%d", i)
		var jsonRepresentation []byte
		var errMashalling error
		var err Error = testCase.err // cast the ErrorForm to Error
		jsonRepresentation, errMashalling = json.Marshal(err.GetMessage())
		assert.Nil(t, errMashalling)
		assert.Equal(t, testCase.expectedStatusCode, err.GetStatusCode())
		assert.Equal(t, testCase.expectedIsError, testCase.err.IsError())
		assert.Equal(t, testCase.expectedJSON, string(jsonRepresentation))
	}
}

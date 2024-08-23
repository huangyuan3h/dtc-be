package errors

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const expected = "{\"message\":\"" + JSONParseError + "\",\"code\":400}"

func TestToString(t *testing.T) {
	httpError := New(JSONParseError, 400)

	jsonStr := httpError.ToString()

	if !cmp.Equal(jsonStr, expected) {
		t.Error("the error response not equal to expected")
	}
}

func TestGatewayResponse(t *testing.T) {
	httpError := New(JSONParseError, 400)
	result, err := httpError.GatewayResponse()
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(result.Body, expected) {
		t.Error("the error response not equal to expected")
	}

}

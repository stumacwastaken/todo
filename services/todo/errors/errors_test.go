package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalError(t *testing.T) {
	err := InternalError()
	assert.Containsf(t, err.Error(), `"message": "internal error"`, "should contain internal error refernece")
}

func TestUnknownError(t *testing.T) {
	err := UnknownError()
	assert.Containsf(t, err.Error(), `"message": "unknown"`, "should contain unknown error reference")
}

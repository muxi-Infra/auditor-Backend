package errorx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSDKError_Error(t *testing.T) {
	var e = MarshalErr(errors.New("test"))
	ok := errors.As(e, &DefaultErr)
	assert.Equal(t, true, ok)
}

package ogitransformer

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	logger "github.com/OpenChaos/ogi/logger"
)

var (
	countOfMockTransformerCalled = 0
	countOfMockTransformerError  = 0
)

type MockTransformer struct{}

func init() {
	logger.SetupLogger()
}

func (m *MockTransformer) Transform(msgid string, msg []byte) ([]byte, error) {
	countOfMockTransformerCalled = 1
	if len(msg) == 0 {
		countOfMockTransformerError += 1
		return []byte{}, errors.New("transform failed")
	}
	return []byte{}, nil
}

func TestTransformSuccess(t *testing.T) {
	var transparentTransformerB bool
	monkey.Patch(NewTransparentTransformer, func() Transformer {
		transparentTransformerB = true
		return &MockTransformer{}
	})

	_, e := Transform("ulid", []byte("{}"))
	assert.Equal(t, e, nil)

	assert.True(t, transparentTransformerB)
	assert.Equal(t, 1, countOfMockTransformerCalled)
	assert.Equal(t, 0, countOfMockTransformerError)
	countOfMockTransformerCalled = 0 // resetting
}

func TestTransformFailure(t *testing.T) {
	var transparentTransformerB bool
	monkey.Patch(NewTransparentTransformer, func() Transformer {
		transparentTransformerB = true
		return &MockTransformer{}
	})

	_, e := Transform("id", []byte(""))
	assert.Equal(t, e.Error(), "transform failed")

	assert.True(t, transparentTransformerB)
	assert.Equal(t, 1, countOfMockTransformerCalled)
	assert.Equal(t, 1, countOfMockTransformerError)
	countOfMockTransformerCalled = 0 // resetting
	countOfMockTransformerError = 0  //resetting
}

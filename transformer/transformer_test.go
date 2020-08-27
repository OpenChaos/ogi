package ogitransformer

import (
	"errors"
	"net/http"
	"testing"

	"bou.ke/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
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

func (m *MockTransformer) Transform(msg []byte) error {
	countOfMockTransformerCalled = 1
	if len(msg) == 0 {
		countOfMockTransformerError += 1
		return errors.New("transform failed")
	}
	return nil
}

func TestTransformSuccess(t *testing.T) {
	var nr, nrEnd *monkey.PatchGuard
	var nrB, nrEndB, transparentTransformerB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		nrB = true
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		nrEndB = true
	})
	monkey.Patch(NewTransparentTransformer, func() Transformer {
		transparentTransformerB = true
		return &MockTransformer{}
	})

	Transform([]byte("{}"))

	assert.True(t, nrB)
	assert.True(t, nrEndB)
	assert.True(t, transparentTransformerB)
	assert.Equal(t, 1, countOfMockTransformerCalled)
	assert.Equal(t, 0, countOfMockTransformerError)
	countOfMockTransformerCalled = 0 // resetting
}

func TestTransformFailure(t *testing.T) {
	var nr, nrEnd *monkey.PatchGuard
	var nrB, nrEndB, transparentTransformerB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		nrB = true
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		nrEndB = true
	})
	monkey.Patch(NewTransparentTransformer, func() Transformer {
		transparentTransformerB = true
		return &MockTransformer{}
	})

	Transform([]byte(""))
	assert.True(t, nrB)
	assert.True(t, nrEndB)
	assert.True(t, transparentTransformerB)
	assert.Equal(t, 1, countOfMockTransformerCalled)
	assert.Equal(t, 1, countOfMockTransformerError)
	countOfMockTransformerCalled = 0 // resetting
	countOfMockTransformerError = 0  //resetting
}

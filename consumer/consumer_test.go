package ogiconsumer

import (
	"net/http"
	"testing"

	"bou.ke/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
)

var (
	countOfMockConsumerCalled = 0
)

type MockConsumer struct{}

func (m *MockConsumer) Consume() {
	countOfMockConsumerCalled += 1
	return
}

func TestConsume(t *testing.T) {
	var nrB, nrEndB bool
	logger.SetupLogger()

	monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nrB = true
		return nil
	})
	monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEndB = true
	})
	monkey.Patch(NewTCPServer, func() Consumer {
		return &MockConsumer{}
	})
	Consume()

	assert.True(t, nrB)
	assert.True(t, nrEndB)
	assert.Equal(t, 1, countOfMockConsumerCalled)
	countOfMockConsumerCalled = 0 //resetting
}

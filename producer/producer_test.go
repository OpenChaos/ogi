package ogiproducer

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
	countOfMockProducerProduceCalled = 0
	countOfMockProducerCloseCalled   = 0
)

func init() {
	logger.SetupLogger()
}

type MockProducer struct{}

func (m *MockProducer) Close() {
	countOfMockProducerCloseCalled += 1
	return
}
func (m *MockProducer) Produce(t string, msg []byte, key string) {
	countOfMockProducerProduceCalled += 1
	return
}

func TestNewProducer(t *testing.T) {
	echoProducer := NewProducer()
	assert.Equal(t, &Echo{}, echoProducer)
}

func TestProduce(t *testing.T) {
	var nrB, nrEndB bool
	mockProducer := &MockProducer{}
	monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nrB = true
		return nil
	})
	monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEndB = true
	})
	monkey.Patch(NewProducer, func() Producer {
		return mockProducer
	})

	Produce("topik", []byte{}, "mykey")
	assert.True(t, nrB)
	assert.True(t, nrEndB)
	assert.Equal(t, countOfMockProducerProduceCalled, 1)
	assert.Equal(t, countOfMockProducerCloseCalled, 1)
	countOfMockProducerProduceCalled = 0 // resetting
	countOfMockProducerCloseCalled = 0   // resetting
}

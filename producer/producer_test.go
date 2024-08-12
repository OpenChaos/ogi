package ogiproducer

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

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
func (m *MockProducer) Produce(msgid string, msg []byte) ([]byte, error) {
	countOfMockProducerProduceCalled += 1
	return []byte{}, nil
}

func TestNewProducer(t *testing.T) {
	echoProducer := NewProducer()
	assert.Equal(t, &Echo{}, echoProducer)
}

func TestProduce(t *testing.T) {
	mockProducer := &MockProducer{}
	monkey.Patch(NewProducer, func() Producer {
		return mockProducer
	})

	_, e := Produce("ulid", []byte{})
	assert.Equal(t, e, nil)

	assert.Equal(t, countOfMockProducerProduceCalled, 1)
	assert.Equal(t, countOfMockProducerCloseCalled, 1)
	countOfMockProducerProduceCalled = 0 // resetting
	countOfMockProducerCloseCalled = 0   // resetting
}

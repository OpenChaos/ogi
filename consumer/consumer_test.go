package ogiconsumer

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

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
	logger.SetupLogger()

	monkey.Patch(NewTCPServer, func() Consumer {
		return &MockConsumer{}
	})
	Consume()

	assert.Equal(t, 1, countOfMockConsumerCalled)
	countOfMockConsumerCalled = 0 //resetting
}

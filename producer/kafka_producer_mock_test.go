package ogiproducer

import (
	logger "github.com/OpenChaos/ogi/logger"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (k *MockProducer) Close() {
	k.Mock.Called()
}

func (k *MockProducer) Produce(topic string, message []byte, messageKey string) {
	k.Mock.Called()
}

func setTestConfig() {
	BootstrapServers = "someserver"
	ProducerType = "confluent-kafka"
	logger.SetupLogger()
}

func unsetTestConfig() {
	ProducerType = ""
	BootstrapServers = ""
}

package ogiproducer

import (
	"github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
	"github.com/gol-gol/golenv"
)

type Producer interface {
	Produce(string, []byte, string)
	Close()
}

type NewProducerFunc func() Producer

var (
	ProducerType = golenv.OverrideIfEnv("PRODUCER_TYPE", "echo")

	producerMap = map[string]NewProducerFunc{
		"echo":   NewEchoProducer,
		"plugin": NewProducerPlugin,
	}
)

func NewProducer() Producer {
	return producerMap[ProducerType]()
}

func Produce(topic string, message []byte, messageKey string) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	producer := NewProducer()
	defer producer.Close()

	producer.Produce(topic, message, messageKey)
	logger.Infof("topic '%s' message-key '%s'", topic, messageKey)
}

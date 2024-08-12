package ogiproducer

import (
	"github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
	"github.com/gol-gol/golenv"
)

type Producer interface {
	Produce(string, []byte) ([]byte, error)
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

func Produce(msgid string, msg []byte) ([]byte, error) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	producer := NewProducer()
	defer producer.Close()

	logger.Infof("msg#[%s]", msgid)
	return producer.Produce(msgid, msg)
}

package ogiconsumer

import (
	"github.com/abhishekkr/gol/golenv"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
)

type Consumer interface {
	Consume()
}

type NewConsumerFunc func() Consumer

var (
	ConsumerType = golenv.OverrideIfEnv("CONSUMER_TYPE", "tcp-server")

	consumerMap = map[string]NewConsumerFunc{
		"tcp-server": NewTCPServer,
		"plugin":     NewConsumerPlugin,
	}
)

func Consume() {
	txn := instrumentation.StartTransaction("consume_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	consumer := consumerMap[ConsumerType]()
	consumer.Consume()
}

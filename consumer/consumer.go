package ogiconsumer

import (
	"github.com/abhishekkr/gol/golenv"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
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

func init() {
	validateConfig()
}

func validateConfig() {
	var missingVariables string

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func failIfError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func Consume() {
	txn := instrumentation.StartTransaction("consume_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	consumer := consumerMap[ConsumerType]()
	consumer.Consume()
}

package ogiconsumer

import (
	"path"
	"plugin"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
	"github.com/abhishekkr/gol/golenv"
)

type ConsumerPlugin struct {
	Name        string
	ConsumeFunc plugin.Symbol
}

var (
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_PLUGIN_PATH", "./consumer.so")
)

func NewConsumerPlugin() Consumer {
	p, err := plugin.Open(ConsumerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin at %s: %s", ConsumerPluginPath, err)
	}

	consumeFunc, err := p.Lookup("Consume")
	if err != nil {
		logger.Fatalf("Error looking up 'Consume': %s", err)
	}

	return &ConsumerPlugin{
		Name:        path.Base(ConsumerPluginPath),
		ConsumeFunc: consumeFunc,
	}
}

func (plugin *ConsumerPlugin) Consume() {
	txn := instrumentation.StartTransaction("event_kafka_message_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	plugin.ConsumeFunc.(func())()
}

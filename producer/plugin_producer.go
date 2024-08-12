package ogiproducer

import (
	"path"
	"plugin"

	logger "github.com/OpenChaos/ogi/logger"
	"github.com/gol-gol/golenv"
)

type ProducerPlugin struct {
	Name        string
	CloseFunc   plugin.Symbol
	ProduceFunc plugin.Symbol
}

var (
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_PLUGIN_PATH", "./producer.so")
)

func NewProducerPlugin() Producer {
	p, err := plugin.Open(ProducerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin: %s", err)
	}

	closeFunc, err := p.Lookup("Close")
	if err != nil {
		logger.Fatalf("Error looking up 'Close': %s", err)
	}

	produceFunc, err := p.Lookup("Produce")
	if err != nil {
		logger.Fatalf("Error looking up 'Produce': %s", err)
	}

	return &ProducerPlugin{
		Name:        path.Base(ProducerPluginPath),
		CloseFunc:   closeFunc,
		ProduceFunc: produceFunc,
	}
}

func (plugin *ProducerPlugin) Close() {
	plugin.CloseFunc.(func())()
}

func (plugin *ProducerPlugin) Produce(topic string, message []byte, messageKey string) {
	plugin.ProduceFunc.(func(string, []byte, string))(topic, message, messageKey)
}

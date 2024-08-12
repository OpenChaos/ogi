package ogitransformer

import (
	"path"
	"plugin"

	"github.com/gol-gol/golenv"

	logger "github.com/OpenChaos/ogi/logger"
)

type TransformerPlugin struct {
	Name          string
	TransformFunc plugin.Symbol
}

var (
	TransformerPluginPath = golenv.OverrideIfEnv("TRANSFORMER_PLUGIN_PATH", "./transformer.so")
)

func NewTransformerPlugin() Transformer {
	p, err := plugin.Open(TransformerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin: %s", err)
	}

	transformFunc, err := p.Lookup("Transform")
	if err != nil {
		logger.Fatalf("Error looking up 'Transform': %s", err)
	}

	return &TransformerPlugin{
		Name:          path.Base(TransformerPluginPath),
		TransformFunc: transformFunc,
	}
}

func (plugin *TransformerPlugin) Transform(msg []byte) error {
	return plugin.TransformFunc.(func([]byte) error)(msg)
}

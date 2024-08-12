package ogitransformer

import (
	"github.com/gol-gol/golenv"
)

type Transformer interface {
	Transform(string, []byte) ([]byte, error)
}

type NewTransformer func() Transformer

var (
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")

	transformerMap = map[string]NewTransformer{
		"transparent": NewTransparentTransformer,
		"plugin":      NewTransformerPlugin,
	}
)

func Transform(msgid string, msg []byte) ([]byte, error) {
	transformer := transformerMap[TransformerType]()
	return transformer.Transform(msgid, msg)
}

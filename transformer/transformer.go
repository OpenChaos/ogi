package ogitransformer

import (
	"github.com/gol-gol/golenv"

	"github.com/OpenChaos/ogi/instrumentation"

	logger "github.com/OpenChaos/ogi/logger"
)

type Transformer interface {
	Transform(string, []byte) error
}

type NewTransformer func() Transformer

var (
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")

	transformerMap = map[string]NewTransformer{
		"transparent": NewTransparentTransformer,
		"plugin":      NewTransformerPlugin,
	}
)

func Transform(msgid string, msg []byte) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	transformer := transformerMap[TransformerType]()
	if err := transformer.Transform(msgid, msg); err != nil {
		// produce to dead-man-talking
		logger.Warn(err)
	}
}

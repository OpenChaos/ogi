package ogitransformer

import (
	"github.com/abhishekkr/gol/golenv"

	"github.com/OpenChaos/ogi/instrumentation"

	logger "github.com/OpenChaos/ogi/logger"
)

type Transformer interface {
	Transform([]byte) error
}

type NewTransformer func() Transformer

var (
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")

	transformerMap = map[string]NewTransformer{
		"transparent": NewTransparentTransformer,
		"plugin":      NewTransformerPlugin,
	}
)

func Transform(msg []byte) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	transformer := transformerMap[TransformerType]()
	if err := transformer.Transform(msg); err != nil {
		// produce to dead-man-talking topic
		logger.Warn(err)
	}
}

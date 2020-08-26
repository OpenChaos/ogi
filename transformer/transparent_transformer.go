package ogitransformer

import ogiproducer "github.com/OpenChaos/ogi/producer"

type TransparentTransformer struct {
	Topic string
}

func (t *TransparentTransformer) Transform(msg []byte) (err error) {
	ogiproducer.Produce(t.Topic, msg, "")
	return nil
}

func NewTransparentTransformer() Transformer {
	return &TransparentTransformer{Topic: ""}
}

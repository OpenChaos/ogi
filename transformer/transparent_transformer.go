package ogitransformer

import ogiproducer "github.com/OpenChaos/ogi/producer"

type TransparentTransformer struct {
}

func (t *TransparentTransformer) Transform(msgid string, msg []byte) (err error) {
	ogiproducer.Produce(msgid, msg)
	return nil
}

func NewTransparentTransformer() Transformer {
	return &TransparentTransformer{}
}

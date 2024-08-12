package ogitransformer

import ogiproducer "github.com/OpenChaos/ogi/producer"

type TransparentTransformer struct {
}

func (t *TransparentTransformer) Transform(msgid string, msg []byte) (response []byte, err error) {
	return ogiproducer.Produce(msgid, msg)
}

func NewTransparentTransformer() Transformer {
	return &TransparentTransformer{}
}

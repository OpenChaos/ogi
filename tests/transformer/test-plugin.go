package main

type TestTransformer struct {
}

var (
	p *TestTransformer
)

func init() {
	p = new(TestTransformer)
}

func (msgLog *TestTransformer) Transform(msg []byte) (err error) {
	return
}

func Transform(msg []byte) error {
	return p.Transform(msg)
}

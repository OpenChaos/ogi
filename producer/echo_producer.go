package ogiproducer

import (
	"fmt"
	"log"
)

type Echo struct {
}

func (e *Echo) Close() {
	fmt.Println("")
}

func (e *Echo) Produce(msgid string, msg []byte) ([]byte, error) {
	if len(msg) != 0 {
		fmt.Println(string(msg))
	} else {
		log.Println("# received blank message @", msgid)
	}
	return []byte{}, nil
}

func NewEchoProducer() Producer {
	return &Echo{}
}

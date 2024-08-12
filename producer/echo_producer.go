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

func (e *Echo) Produce(msgid string, msg []byte) {
	if len(msg) != 0 {
		fmt.Println(string(msg))
	} else {
		log.Println("# received blank message @", msgid)
	}
}

func NewEchoProducer() Producer {
	return &Echo{}
}

package ogiproducer

import (
	"fmt"
	"log"
	"time"
)

type Echo struct {
}

func (e *Echo) Close() {
	fmt.Println("")
}

func (e *Echo) Produce(topic string, message []byte, messageKey string) {
	if topic != "" {
		fmt.Println("topic:", topic)
	}
	if messageKey != "" {
		fmt.Println("key:", messageKey)
	} else {
		fmt.Println("key:", time.Now())
	}
	if len(message) != 0 {
		fmt.Println(string(message))
	} else {
		log.Println("# received blank message")
	}
}

func NewEchoProducer() Producer {
	return &Echo{}
}

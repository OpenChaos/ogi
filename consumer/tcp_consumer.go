package ogiconsumer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gol-gol/golenv"
	ulid "github.com/oklog/ulid/v2"

	ogitransformer "github.com/OpenChaos/ogi/transformer"
)

type TCPServer struct {
	Port string
}

var (
	tcpConsumerPort = golenv.OverrideIfEnv("OGI_TCP_CONSUMER_PORT", "8080")
)

func (t *TCPServer) Transform(msgid, lyne string) {
	msgB, err := ogitransformer.Transform(msgid, []byte(lyne))
	if err != nil {
		log.Println("[ERROR]", err.Error())
	} else if len(msgB) > 0 {
		log.Println("[RESPONSE]", string(msgB))
	}
}

func (t *TCPServer) Start() {
	l, err := net.Listen("tcp", t.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(netData) == "exit" {
			fmt.Println("bye!")
			return
		}

		msgid := ulid.Make().String()
		t.Transform(msgid, netData)
		tym := time.Now()
		response := msgid + ": " + tym.Format(time.RFC3339) + "\n"
		c.Write([]byte(response))
	}
}

func (t *TCPServer) Consume() {
	log.Printf("tcp server listening at: %s", t.Port)
	t.Start()
}

func NewTCPServer() Consumer {
	tcpServer := TCPServer{Port: ":" + tcpConsumerPort}
	return &tcpServer
}

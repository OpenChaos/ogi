package ogiconsumer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gol-gol/golenv"

	ogitransformer "github.com/OpenChaos/ogi/transformer"
)

type TCPServer struct {
	Port string
}

var (
	tcpConsumerPort = golenv.OverrideIfEnv("OGI_TCP_CONSUMER_PORT", "8080")
)

func (t *TCPServer) Transform(lyne string) {
	ogitransformer.Transform([]byte(lyne))
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
		if strings.TrimSpace(string(netData)) == "exit" {
			fmt.Println("bye!")
			return
		}

		t.Transform(string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
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

package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

type ConnectionMessages struct {
	messages [][]byte
	mutex    sync.RWMutex
}

func main() {
	godotenv.Load()

	var port uint = 0
	flag.UintVar(&port, "port", 0, "the port to run on for the websocket")

	if port == 0 {
		portStr := os.Getenv("PORT")
		portEnv, err := strconv.Atoi(portStr)
		if err == nil {
			port = uint(portEnv)
		}
	}

	assert.Assert(port != 0, "please provide a port for the relay server")

	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

	slog.Warn("port selected", "port", port)
	r := relay.NewRelay(uint16(port), uuid)

	connMsgs := &ConnectionMessages{
		messages: make([][]byte, 0),
		mutex:    sync.RWMutex{},
	}

	go newConnections(r, connMsgs)
	go onMessage(r, connMsgs)

	r.Start()
}

func newConnections(relay *relay.Relay, msgs *ConnectionMessages) {
	for {
		conn := <-relay.NewConnections()
		msgs.mutex.RLock()
		for _, msg := range msgs.messages {
			conn.Conn.WriteMessage(websocket.BinaryMessage, msg)
		}
		msgs.mutex.RUnlock()
	}
}

func onMessage(relay *relay.Relay, msgs *ConnectionMessages) {
	framer := net.NewByteFramer()
	go framer.FrameChan(relay.Messages())
	for {
		frame := <-framer.Frames()
		switch frame.CmdType {
		case byte(net.OPEN), byte(net.BRIGHTNESS_TO_ASCII):
			length := net.HEADER_SIZE + len(frame.Data)
			encoded := make([]byte, length, length)
			frame.Into(encoded, 0)
			msgs.mutex.Lock()
			msgs.messages = append(msgs.messages, encoded)
			msgs.mutex.Unlock()
		}
	}
}

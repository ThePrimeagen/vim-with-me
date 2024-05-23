package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

type ConnectionMessages struct {
	open  []byte
	mutex sync.RWMutex
}

func main() {
	godotenv.Load()

	var port uint = 0
	flag.UintVar(&port, "port", 0, "the port to run on for the websocket")
	var staticMetricsFilename string
	flag.StringVar(&staticMetricsFilename, "metrics-filename", "", "a filename to write the metrics to (plain text)")
	flag.Parse()

	if port == 0 {
		portStr := os.Getenv("PORT")
		portEnv, err := strconv.Atoi(portStr)
		if err == nil {
			port = uint(portEnv)
		}
	}

	if staticMetricsFilename == "" {
		staticMetricsFilename = os.Getenv("METRICS_FILENAME")
	}

	assert.Assert(port != 0, "please provide a port for the relay server")

	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

	slog.Warn("port selected", "port", port)
	r := relay.NewRelay(uint16(port), uuid)

	connMsgs := &ConnectionMessages{
		open:  nil,
		mutex: sync.RWMutex{},
	}

	go newConnections(r, connMsgs)
	go onMessage(r, connMsgs)

	if staticMetricsFilename != "" {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		go func() {
			for {
				<-ticker.C
				current, maxConcurrent, total := r.GetConnectionStats()
				writeErr := os.WriteFile(staticMetricsFilename,
					[]byte(fmt.Sprintf(
						"         Current: %d\n"+
						"  Max concurrent: %d\n"+
						"Cumulatice total: %d",
						current, maxConcurrent, total)),
					0644)
				if writeErr != nil {
					slog.Error(writeErr.Error(), "filename", staticMetricsFilename)
					return
				}
			}
		}()
	}

	r.Start()
}

func newConnections(relay *relay.Relay, msgs *ConnectionMessages) {
	for {
		conn := <-relay.NewConnections()
		if msgs.open != nil {
			msgs.mutex.RLock()
			slog.Warn("new connection, appending open messages", "open", msgs.open)
			conn.Conn.WriteMessage(websocket.BinaryMessage, msgs.open)
			msgs.mutex.RUnlock()
		}
	}
}

func onMessage(relay *relay.Relay, msgs *ConnectionMessages) {
	framer := net.NewByteFramer()
	go framer.FrameChan(relay.Messages())
	for frame := range framer.Frames() {
		slog.Debug("received frame", "frame", frame)

		switch frame.CmdType {
		case byte(net.OPEN):
            bytes := frame.Bytes()
			slog.Warn("new open command", "encoded", bytes)

			msgs.mutex.Lock()
			msgs.open = bytes
			msgs.mutex.Unlock()
		}
	}
}

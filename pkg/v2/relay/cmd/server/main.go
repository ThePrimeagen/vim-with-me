package main

import (
	"flag"
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
	"github.com/theprimeagen/vim-with-me/pkg/v2/metrics"
	"runtime"
)

type ConnectionMessages struct {
	open  []byte
	mutex sync.RWMutex
}

func main() {
	runtime.GOMAXPROCS(runtime.GOMAXPROCS(0) - 1)

	godotenv.Load()

	var port uint = 0
	var staticMetricsFilename, jsonMetricsFilename string
	flag.UintVar(&port, "port", 0, "the port to run on for the websocket")
	flag.StringVar(&staticMetricsFilename, "text-metrics-filename", os.Getenv("STATIC_METRICS_FILENAME"),
		"a filename to periodically update with fresh metrics in text format (every 10 seconds)")
	flag.StringVar(&jsonMetricsFilename, "json-metrics-filename", os.Getenv("JSON_METRICS_FILENAME"),
		"a filename to periodically append JSON metrics to (every 10 seconds)")
	flag.Parse()

	if port == 0 {
		portStr := os.Getenv("PORT")
		portEnv, err := strconv.Atoi(portStr)
		if err == nil {
			port = uint(portEnv)
		}
	}
	assert.Assert(port != 0, "please provide a port for the relay server")

	stats := metrics.New()
	if staticMetricsFilename != "" {
		stats.WithFileWriter(staticMetricsFilename, metrics.FileWriterFormatText, 10 * time.Second)
	}
	if jsonMetricsFilename != "" {
		stats.WithFileWriter(jsonMetricsFilename, metrics.FileWriterFormatAppendJSON, 10 * time.Second)
	}

	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

	slog.Warn("port selected", "port", port)
	r := relay.NewRelay(uint16(port), uuid, stats)

	connMsgs := &ConnectionMessages{
		open:  nil,
		mutex: sync.RWMutex{},
	}

	go newConnections(r, connMsgs)
	go onMessage(r, connMsgs)

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

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/imkira/go-observer/v2"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func toSSE(evt string, bytes []byte) []byte {
	return []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", evt, string(bytes)))
}

func toChannel(conn *tcp.Connection, ctx context.Context) chan *tcp.TCPCommand {
	ch := make(chan *tcp.TCPCommand, 10)

	go func() {
	outer:
		for {
			select {
			case <-ctx.Done():
				break outer
			default:
			}

			msg, err := conn.Next()
			if err != nil {
				slog.Error("error while reading from channel", "err", err)
				break outer
			}

			ch <- msg
		}
	}()

	return ch
}

// 2. CSS grid stuff

func runListener(prop observer.Property[[]byte], ctx context.Context, conn *tcp.Connection, renderer **window.Renderer) error {

	connCh := toChannel(conn, ctx)

	for msg := range connCh {
		if msg.Command == commands.OPEN_WINDOW {
			*renderer = window.NewRender(int(msg.Data[0]), int(msg.Data[1]))
		} else if msg.Command == commands.PARTIAL_RENDER {
			cells, err := commands.PartialRendersFromTCPCommand(msg)
			assert.Assert(err == nil, "partial render failed to parse")
			assert.Assert(cells != nil, "cells is nil")

			(*renderer).FromRemoteRenderer(cells)
			(*renderer).Render()
			fmt.Printf("RENDER STATE: \n%s\n", (*renderer).Debug())
		}

		bytes, err := commands.Jsonify(msg)
		if err != nil {
			return err
		}
		prop.Update(toSSE(commands.CommandNameLookup[msg.Command], bytes))
	}

	return nil
}

type OpenCommand struct {
	Rows int
	Cols int
}

type Commands map[string]int

func sendCmd(w io.Writer, cmd *tcp.TCPCommand) {
	jsonData, err := commands.Jsonify(cmd)
	assert.Assert(err == nil, "should never fail marshaling render command")

	w.Write(toSSE(commands.CommandNameLookup[cmd.Command], jsonData))
	w.(http.Flusher).Flush()
}

func main() {
	var port uint = 0
	host := ""

	flag.StringVar(&host, "host", "127.0.0.1", "host to connect to")
	flag.UintVar(&port, "port", 42069, "port to connect to")
	flag.Parse()

	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		slog.Error("error connecting client", "error", err)
		return
	}
	conn := tcp.NewConnection(c, 0)
	defer conn.Close()

	var innerRenderer *window.Renderer = nil
	var renderer **window.Renderer = &innerRenderer

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	prop := observer.NewProperty([]byte{})

	go func() {
		err := runListener(prop, ctx, &conn, renderer)
		slog.Warn("server finished", "error", err)
		assert.Assert(err == nil, "we somehow errored :(")
	}()

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		stream := prop.Observe()

		// Set CORS headers to allow all origins. You may want to restrict this
		// to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		assert.Assert(*renderer != nil, "SSE: renderer should always be defined before first connection")
		sendCmd(w, commands.OpenCommand((*renderer)))
		sendCmd(w, commands.PartialRender((*renderer).FullRender()))
		fmt.Printf("full render: %s\n", (*renderer).Debug())

	outer:
		for {
			select {
			case <-ctx.Done():
				break outer
			case <-stream.Changes():
				stream.Next()
				val := stream.Value()
				assert.Assert(val != nil, "i should never receive a nil command")

				fmt.Printf("render: %d\n", len(val))
				if len(val) == 0 {
					continue
				}

				w.Write(val)
				w.(http.Flusher).Flush()
			}
		}

	})

	http.ListenAndServe(":8080", nil)
}

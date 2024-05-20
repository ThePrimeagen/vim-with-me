package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

type Opts struct {
	port        uint
	host        string
	parallel    int
	count       int
	connections int
}

type Stats struct {
	totalBytes    int
	totalCommands int

	mutex sync.Mutex
}

func (s *Stats) add(cmds, bytes int) {
	s.mutex.Lock()
	s.totalCommands += cmds
	s.totalBytes += bytes
	s.mutex.Unlock()
}

func parseOpts() Opts {

	opts := Opts{}

	flag.UintVar(&opts.port, "port", 42069, "Port to listen on")
	flag.StringVar(&opts.host, "host", "127.0.0.1", "Host to listen on")
	flag.IntVar(&opts.parallel, "parallel", 1000, "the amount of parallel connections to the server")
	flag.IntVar(&opts.count, "count", 100_000, "the number of requests to receive before disconnecting")
	flag.IntVar(&opts.connections, "connections", 1000, "how many connections to make")

	flag.UintVar(&opts.port, "p", 42069, "Port to listen on")
	flag.StringVar(&opts.host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&opts.parallel, "q", 1_000, "the amount of parallel connections to the server")
	flag.IntVar(&opts.count, "c", 1_000, "the number of requests to receive before disconnecting")
	flag.IntVar(&opts.connections, "x", 1000, "how many connections to make")

	flag.Parse()
	return opts
}

func connect(opts *Opts, stats *Stats, semaphore chan struct{}, id int) int {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.host, opts.port))
	if err != nil {
		slog.Warn("error connecting client", "error", err)
		semaphore <- struct{}{}
		return -1
	}

	totalN := 0
	totalCommands := 0
	_ = id
	defer c.Close()

	conn := tcp.NewConnection(c, id)

	for count := opts.count; count > 0; count-- {
		cmd, err := conn.Next()
		totalCommands += 1
		if err != nil {
			slog.Warn("error connecting client", "error", err)
			if !errors.Is(err, io.EOF) {
			}
			break
		}
		totalN += len(cmd.Data) + tcp.HEADER_SIZE
	}

	semaphore <- struct{}{}

	stats.add(totalCommands, totalN)

	return totalN
}

func main() {

	testies.SetupLogger()

	opts := parseOpts()
	fmt.Printf("%+v\n", &opts)

	stats := Stats{
		totalBytes:    0,
		totalCommands: 0,
		mutex:         sync.Mutex{},
	}

	semaphore := make(chan struct{}, opts.parallel)
	for i := 0; i < opts.parallel && i < opts.connections; i++ {
		semaphore <- struct{}{}
	}

	id := 0
	for connections := opts.connections; connections > 0; connections-- {
		<-semaphore
		id++
		go connect(&opts, &stats, semaphore, id)
	}

	for i := 0; i < opts.parallel && i < opts.connections; i++ {
		<-semaphore
	}

	fmt.Printf("%+v\n", &stats)
}

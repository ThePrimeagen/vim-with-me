package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"

	//"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	"github.com/theprimeagen/vim-with-me/pkg/v2/program"
)

type RelayClient struct {
	client *relay.RelayDriver
	cache  []byte
}

func NewRelayClient(r string) (*RelayClient, error) {
	if len(r) == 0 {
		return &RelayClient{}, nil
	}

	uuid := os.Getenv("AUTH_ID")
	length := 256 * 256
	client := &RelayClient{
		client: relay.NewRelayDriver(r, "/ws", uuid),
		cache:  make([]byte, length, length),
	}

	return client, client.client.Connect()
}

func (r *RelayClient) send(frameable *net.Frameable) {
	n, err := frameable.Into(r.cache, 0)

	assert.NoError(err, "relay server could not call frame#into")

	err = r.client.Relay(r.cache[:n])
	assert.NoError(err, "relay client errored")
}

func (r *RelayClient) sendFrame(frame *encoder.EncodingFrame) {
	if r.client == nil {
		return
	}

	fmt.Printf("sending frame into relay(%s): %d\n", encoder.EncodingName(frame.Encoding), frame.Len)

	r.send(&net.Frameable{Item: frame})
}

func compareFrame(data []byte, encFrame *encoder.EncodingFrame) {
    assert.Assert(false, "unimplemented")
}

func main() {
	godotenv.Load()

	debug := ""
	flag.StringVar(&debug, "debug", "", "runs the file like the program instead of running doom")

	assertF := ""
	flag.StringVar(&assertF, "assert", "", "add an assert file")

	compare := false
	flag.BoolVar(&compare, "compare", false, "compare the encoded and decoded values")

	rounds := 1000000
	flag.IntVar(&rounds, "rounds", 1000000, "the rounds of doom to play")

	relayStr := ""
	flag.StringVar(&relayStr, "relay", "", "the relay server to attach to")
	flag.Parse()

	args := flag.Args()
	name := args[0]

	fmt.Printf("assert file attached \"%s\"\n", assertF)
	fmt.Printf("debug file attached \"%s\"\n", debug)
	fmt.Printf("args file attached \"%v\"\n", args)
	fmt.Printf("relay \"%v\"\n", relayStr)

	relay, err := NewRelayClient(relayStr)
	assert.NoError(err, "failed attempting to connect to server")

	d := doom.NewDoom()

	prog := program.
		NewProgram(name).
		WithArgs(args[1:]).
		WithWriter(d)

	if debug != "" {
		debugFile, err := os.Create(debug)
		assert.NoError(err, "unable to open debug file")
		prog = prog.WithWriter(debugFile)
	}

	if assertF != "" {
		assertFile, err := os.Create(assertF)
		assert.NoError(err, "unable to open assert file")
		assert.ToWriter(assertFile)
	}

	ctx := context.Background()
	go func() {
		err := prog.Run(ctx)
		assert.NoError(err, "prog.Run(ctx)")
	}()

	<-d.Ready()

	enc := encoder.NewEncoder(d.Rows*(d.Cols/2), ascii_buffer.QuadtreeParam{
		Depth:  2,
		Stride: 1,
		Rows:   d.Rows,
		Cols:   d.Cols / 2,
	})

	enc.AddEncoder(encoder.XorRLE)
	enc.AddEncoder(encoder.Huffman)

	relay.send(net.CreateOpen(d.Rows, d.Cols))

	frames := d.Frames()

	go func() {
		<-time.After(time.Second * 1)
		prog.SendKey('')
		<-time.After(time.Second * 1)
		prog.SendKey('')
		<-time.After(time.Second * 1)
		prog.SendKey('')
		<-time.After(time.Second * 1)
		prog.SendKey('')
		<-time.After(time.Second * 1)
		prog.SendKey('')

		count := 1
		for {
			<-time.After(time.Millisecond * 16)
			if count%100 == 0 {
				prog.SendKey('f')
			} else {
				prog.SendKey('w')
			}
			count++

		}
	}()

	for range rounds {
		select {
		case frame := <-frames:
			data := ansiparser.RemoveAsciiStyledPixels(frame.Color)

			encFrame := enc.PushFrame(data)
			assert.NotNil(encFrame, "expected enc frame to be not nil")
			relay.sendFrame(encFrame)

			if compare {
                compareFrame(data, encFrame)
			}
		}
	}
}

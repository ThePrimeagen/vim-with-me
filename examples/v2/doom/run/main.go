package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"

	//"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	"github.com/theprimeagen/vim-with-me/pkg/v2/program"
)

type RelayClient struct {
	client *relay.RelayDriver
	cache  []byte
	url    string
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
		url:    r,
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

	r.send(&net.Frameable{Item: frame})
}

func compareFrame(data []byte, encFrame *encoder.EncodingFrame) {
	assert.Assert(false, "unimplemented")
}

func runChat(prog *program.Program) {
	ctx := context.Background()

	//doom create controller
	twitchChat, err := chat.NewTwitchChat(ctx, "theprimeagen")
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(doom.DoomFilterFn).
		WithMap(doom.DoomChatMapFn)
	go chtAgg.Pipe(twitchChat)
	doomCtrl := doom.NewDoomController(prog)
	ctrl := controller.
		NewController(&chtAgg, doomCtrl).
		WithInputTimer(time.NewTicker(time.Millisecond * 250).C).
		WithPlayTimer(time.NewTicker(time.Millisecond * 20).C)
	go ctrl.Start(ctx)

	go func() {
		<-time.After(time.Millisecond * 500)
		prog.SendKey("")
		<-time.After(time.Millisecond * 500)
		prog.SendKey("")
		<-time.After(time.Millisecond * 500)
		prog.SendKey("")
		<-time.After(time.Millisecond * 500)
		prog.SendKey("")
		<-time.After(time.Millisecond * 500)
		prog.SendKey("")

		doomCtrl.Play()
	}()
}

func send(relays []*RelayClient, frameable *net.Frameable) {
	for _, relay := range relays {
		relay.send(frameable)
	}
}

func createEnc(encoderFn encoder.EncodingCall, rows, cols int) *encoder.Encoder {
	enc := encoder.NewEncoder(rows*(cols/2), ascii_buffer.QuadtreeParam{
		Depth:  2,
		Stride: 1,
		Rows:   rows,
		Cols:   cols / 2,
	})

	enc.AddEncoder(encoderFn)
	return enc
}

type EncodingStats struct {
	Unencoded               int64
	HalfEncoding            int64
	HalfAndNoChars          int64
	HalfNoCharsReducedColor int64
	DropDuplicateFrames     int64
	RLE                     int64
	XorRLE                  int64
	Huffman                 int64
	Combined                int64
	Frame                   int64
}

func (e *EncodingStats) Header() string {
	return "Frame,Unencoded,HalfEncoding,HalfAndNoChars,HalfNoCharsReducedColor,DropDuplicateFrames,RLE,XorRLE,Huffman,Combined"
}

func (e *EncodingStats) Stat() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d",
		e.Frame,
		e.Unencoded / 1000,
		e.HalfEncoding / 1000,
		e.HalfAndNoChars / 1000,
        e.HalfNoCharsReducedColor / 1000,
		e.DropDuplicateFrames / 1000,
		e.RLE / 1000,
		e.XorRLE / 1000,
		e.Huffman / 1000,
		e.Combined / 1000,)
}

func main() {
	godotenv.Load()

	debug := ""
	flag.StringVar(&debug, "debug", "", "runs the file like the program instead of running doom")

	assertF := ""
	flag.StringVar(&assertF, "assert", "", "add an assert file")

	compare := false
	flag.BoolVar(&compare, "compare", false, "compare the encoded and decoded values")

	allowChat := false
	flag.BoolVar(&allowChat, "chat", false, "allow for chat interfacing")

	displayOutput := false
	flag.BoolVar(&displayOutput, "display", false, "displays doom in the terminal")

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

	relays := make([]*RelayClient, 0)
	relayDests := strings.Split(relayStr, ",")

	for _, r := range relayDests {
		if r == "" {
			continue
		}
		relay, err := NewRelayClient(r)
		assert.NoError(err, "failed attempting to connect to server")

		relays = append(relays, relay)
	}

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

	none := createEnc(encoder.NoneEncoding, d.Rows, d.Cols)
	rle := createEnc(encoder.RLEEncoding, d.Rows, d.Cols).AddEncoder(encoder.NoneEncoding)
	rleXor := createEnc(encoder.XorRLE, d.Rows, d.Cols).AddEncoder(encoder.NoneEncoding)
	huff := createEnc(encoder.Huffman, d.Rows, d.Cols).AddEncoder(encoder.NoneEncoding)

	enc := createEnc(encoder.NoneEncoding, d.Rows, d.Cols).
		AddEncoder(encoder.RLEEncoding).
		AddEncoder(encoder.XorRLE).
		AddEncoder(encoder.Huffman)
	encStats := &EncodingStats{}

	send(relays, net.CreateOpen(d.Rows, d.Cols))

	frames := d.Frames()

	if allowChat {
		runChat(prog)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	fmt.Printf("doom: %d x %d\n", d.Rows, d.Cols)
	fmt.Printf("%s\n", encStats.Header())
	go func() {
		<-c
		fmt.Printf("%s\n", encStats.Stat())
		os.Exit(1)
	}()

	for range rounds {
		if encStats.Frame%250 == 0 {
			fmt.Printf("%s\n", encStats.Stat())
		}
		encStats.Frame++
		select {
		case frame := <-frames:
			data := ansiparser.RemoveAsciiStyledPixels(frame.Color)
			bytes := int64(d.Rows) * int64(d.Cols)
			encStats.Unencoded += bytes * 4 // 3 bytes color, 1 char
			encStats.HalfEncoding += (bytes * 4) / 2
			encStats.HalfAndNoChars += (bytes * 3) / 2
			encStats.HalfNoCharsReducedColor += bytes / 2

			encFrame := enc.PushFrame(data)
			noneFrame := none.PushFrame(data)
			rleFrame := rle.PushFrame(data)
			rleXorFrame := rleXor.PushFrame(data)
			huffFrame := huff.PushFrame(data)

			if encFrame != nil {
				encStats.Combined += int64(encFrame.Len)
			}

			if noneFrame != nil {
				encStats.DropDuplicateFrames += int64(noneFrame.Len)
			}

			if rleFrame != nil {
				encStats.RLE += int64(rleFrame.Len)
			}

			if rleXorFrame != nil {
				encStats.XorRLE += int64(rleXorFrame.Len)
			}

			if huffFrame != nil {
				encStats.Huffman += int64(huffFrame.Len)
			}

			if encFrame == nil {
				break
			}

			send(relays, net.NewFrameable(encFrame))

			if compare {
				compareFrame(data, encFrame)
			}

			if displayOutput {
				fmt.Printf(display.Display(&frame, d.Rows, d.Cols))
			}
		}
	}

}

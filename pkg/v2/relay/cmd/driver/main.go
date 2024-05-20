package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

func main() {
	err := godotenv.Load()
	assert.NoError(err, "unable to load dotenv")

	var addr string = ""
	flag.StringVar(&addr, "addr", "", "the address for the relay driver to send messages to")
	addr = "vim-with-me.fly.dev:8080"

	assert.Assert(addr != "", "driver requires an addr")

	var f string = ""
	flag.StringVar(&f, "file", "/tmp/out", "the file of the client driver")

	file, err := os.Open(f)
	assert.NoError(err, "driver file error")

	fmt.Printf("relay: %s\n", addr)
	client := relay.NewRelayDriver(addr, os.Getenv("AUTH_ID"))
	err = client.Connect()
	assert.NoError(err, "unable to connect to relay")
	defer client.Close()

	lines := bufio.NewScanner(file)
	for lines.Scan() {
		txt := lines.Text()
		fmt.Printf("sending line: %s\n", txt)
		err := client.Relay([]byte(txt))
		assert.NoError(err, "unable to relay data")

		<-time.NewTimer(time.Second).C
	}
}

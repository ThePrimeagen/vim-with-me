package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

type Msg struct {
    msg []byte
    idx int
    count int
}

func client(port uint16, idx int, out chan<- Msg) {
    count := 0

    u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%d", port), Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    defer c.Close()
    assert.NoError(err, "unable to connect to the websocket server")

    for {
        _, msg, err := c.ReadMessage()
        assert.NoError(err, "read message failed")

        out <- Msg{
            count: count,
            idx: idx,
            msg: msg,
        }

        count++
    }
}

func main() {
    godotenv.Load()
	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

    port := uint16(42069)

    fmt.Printf("creating relay\n")
    r := relay.NewRelay(port, uuid)
    go r.Start()

    <-time.NewTimer(time.Millisecond * 500).C
    <-time.NewTimer(time.Millisecond * 500).C

    ch := make(chan Msg, 250)

    go client(port, 1, ch)
    go client(port, 2, ch)
    go client(port, 3, ch)

    fmt.Printf("created driver\n")
    client := relay.NewRelayDriver(fmt.Sprintf("localhost:%d", port), os.Getenv("AUTH_ID"))
    err := client.Connect()
    assert.NoError(err, "unable to connect to relay")
    defer client.Close()
    <-time.NewTimer(time.Millisecond * 500).C

    line := []string{
        "aoeu",
        "aoeuaoeu",
        "aoeuaoeuaoeu",
        "aoeuaoeu",
        "aoeu",
    }

    for i, l := range line {
        fmt.Printf("for %s\n", l)
        err := client.Relay([]byte(l))
        assert.NoError(err, "unable to relay data")

        for range 3 {
            select {
            case msg := <-ch:
                assert.Assert(i == msg.count, "expecting message #count to equal i", "count", msg.count, "i", i)
                assert.Assert(string(msg.msg) == l, "expecting msg == line", "msg.msg", string(msg.msg), "line", l)
            case <-time.NewTimer(time.Second).C:
                assert.Assert(false, "waiting for message", "line", l)
            }
        }
    }


}

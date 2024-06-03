package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var count = 0
func makeMessage() []byte {
    count++

    size := 250
    if count & 0x1 == 1 {
        size = 5000
    }

    out := make([]byte, size, size)
    out[0] = byte(count)

    return out
}

var upgrader = websocket.Upgrader{} // use default options
func main() {

    upgrader.CheckOrigin = func(r *http.Request) bool {
        return true
    }

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        c, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Fatal("upgrade", err)
        }

        for {
            <-time.After(time.Millisecond * 16)
            c.WriteMessage(websocket.BinaryMessage, makeMessage())
        }
	})

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	slog.Warn("listening and serving http", "http", addr)
	err := http.ListenAndServe(addr, nil)

	log.Fatal(err)

}

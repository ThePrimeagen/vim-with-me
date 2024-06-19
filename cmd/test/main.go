package main


func Backward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s)-1; i >= 0; i-- {
			if !yield(i, s[i]) {
				// Where clean-up code goes
				return
			}
		}
	}
}

// Explicit struct for the state
Backward_Iterator :: struct($E: typeid) {
	slice: []E,
	idx:   int,
}

// Explicit construction for the iterator
backward_make :: proc(s: []$E) -> Backward_Iterator(E) {
	return {slice = s, idx = len(s)-1}
}

backward_iterate :: proc(it: ^Backward_Iterator($E)) -> (elem: E, idx: int, ok: bool) {
	if it.idx >= 0 {
		elem, idx, ok = it.slice[it.idx], it.idx, true
		it.idx -= 1
	}
	return
}


























import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Foo struct {
    ...
}



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

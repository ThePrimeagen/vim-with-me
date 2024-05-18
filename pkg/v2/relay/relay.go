package relay

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type Relay struct {
	WSPort uint16
	uuid   string

	mutex     sync.RWMutex
	listeners map[int]*Conn

	id int
}

var upgrader = websocket.Upgrader{} // use default options

func NewRelay(ws uint16) *Relay {
	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

	return &Relay{
		WSPort: ws,
		uuid:   uuid,
	}
}

func (relay *Relay) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		relay.render(w, r)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", relay.WSPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (relay *Relay) relay(data []byte) {
	relay.mutex.RLock()
	for _, conn := range relay.listeners {
        conn.msg(data)
	}
	relay.mutex.RUnlock()
}

func (relay *Relay) remove(id int) {
	relay.mutex.Lock()
	delete(relay.listeners, id)
	relay.mutex.Unlock()
}

func (relay *Relay) add(id int, ws *websocket.Conn) {
    conn := &Conn{
        id: id,
        conn: ws,
        msgs: make(chan []byte, 10),
    }

	relay.mutex.Lock()
	relay.listeners[id] = conn
	relay.mutex.Unlock()

    go conn.read()
    go conn.write()
}

func (relay *Relay) render(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	id := relay.id
    relay.add(id, c)

	relay.id++

}

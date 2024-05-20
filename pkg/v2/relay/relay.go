package relay

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Relay struct {
	port uint16
	uuid string

	ch chan []byte

	mutex     sync.RWMutex
	listeners map[int]*Conn

	id int
}

var upgrader = websocket.Upgrader{} // use default options

func NewRelay(ws uint16, uuid string) *Relay {

	return &Relay{
		port: ws,
		uuid: uuid,

		ch: make(chan []byte, 10),

		mutex:     sync.RWMutex{},
		listeners: make(map[int]*Conn),

		id: 0,
	}
}

func (relay *Relay) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        slog.Warn("new ws command coming in")
		relay.render(w, r)
	})

    count := 0
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
        count++
        slog.Warn("healthcheck", "count", count)
    })

	addr := fmt.Sprintf("127.0.0.1:%d", relay.port)
    slog.Warn("listening and serving http", "http", addr)
    err := http.ListenAndServe(addr, nil)

    log.Fatal(err)
}

func (relay *Relay) Messages() chan []byte {
	return relay.ch
}

func (relay *Relay) relay(data []byte) {
    // quick write to prevent blocking if there is no listener
    select {
    case relay.ch <- data:
    default:
    }
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
    slog.Warn("removing client connection", "id", id)
}

func (relay *Relay) add(id int, ws *websocket.Conn) {
	conn := &Conn{
		id:    id,
		conn:  ws,
		msgs:  make(chan []byte, 10),
		relay: relay,
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
    slog.Warn("connection established", "id", id)

	relay.id++
}

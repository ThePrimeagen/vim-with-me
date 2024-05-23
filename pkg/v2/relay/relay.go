package relay

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Relay struct {
	port uint16
	uuid string

	ch    chan []byte
	conns chan *Conn

	maxConcurrentConnections int
	totalConnections         int

	mutex     sync.RWMutex
	listeners map[int]*Conn

	id int
}

var upgrader = websocket.Upgrader{} // use default options

func NewRelay(ws uint16, uuid string) *Relay {

	return &Relay{
		port: ws,
		uuid: uuid,

		ch:    make(chan []byte, 10),
		conns: make(chan *Conn, 10),

		mutex:     sync.RWMutex{},
		listeners: make(map[int]*Conn),

		id: 0,
	}
}

func (relay *Relay) GetConnectionStats() (int, int, int) {
    relay.mutex.RLock()
    defer relay.mutex.RUnlock()
    return len(relay.listeners), relay.maxConcurrentConnections, relay.totalConnections
}

//TODO(post doom): Fix this shit
/** THIS IS SHITTY **/
/** Relay should really just be something in which i hand connections and that
* is it, no concept here.  Maybe not even connections but writers */
/** THIS IS SHITTY /end **/
func (relay *Relay) Start() {
	tmpl := template.Must(template.ParseGlob("./views/*.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", struct{}{})
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		current, maxConcurrent, total := relay.GetConnectionStats()
		w.Write([]byte("# help connections_current total number of current connections\n"))
		w.Write([]byte("# type connections_current gauge\n"))
		w.Write([]byte("connections_current_total " + strconv.Itoa(current) + "\n"))
		w.Write([]byte("# help connections_maxconcurrent maximum number of concurrent connections\n")) 
		w.Write([]byte("# type connections_maxconcurrent gauge\n"))
		w.Write([]byte("connections_maxconcurrent_total " + strconv.Itoa(maxConcurrent) + "\n"))
		w.Write([]byte("# help connections_total number of connections\n"))
		w.Write([]byte("# type connections_total counter\n"))
		w.Write([]byte("connections_total " + strconv.Itoa(total) + "\n"))
	})

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		relay.render(w, r)
	})

	count := 0
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		count++
		slog.Warn("healthcheck", "count", count)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", relay.port)
	slog.Warn("listening and serving http", "http", addr)
	err := http.ListenAndServe(addr, nil)

	log.Fatal(err)
}

func (relay *Relay) Messages() chan []byte {
	return relay.ch
}

func (relay *Relay) NewConnections() chan *Conn {
	return relay.conns
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
		Conn:  ws,
		msgs:  make(chan []byte, 10),
		relay: relay,
	}

	relay.mutex.Lock()
    
    relay.listeners[id] = conn

    relay.totalConnections++
    numConnections := len(relay.listeners)
    if numConnections > relay.maxConcurrentConnections {
        relay.maxConcurrentConnections = numConnections
    }
	
    relay.mutex.Unlock()

	select {
	case relay.conns <- conn:
	default:
	}

	go conn.read()
	go conn.write()
}

func (relay *Relay) render(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// force sync on ids
	// probably shoud look into atomic ints here...
	relay.mutex.Lock()
	relay.id++
	id := relay.id
	relay.mutex.Unlock()

	relay.add(id, c)
	slog.Warn("connection established", "id", id)

}

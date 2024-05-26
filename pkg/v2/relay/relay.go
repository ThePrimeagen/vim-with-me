package relay

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/metrics"
)

type Relay struct {
	port uint16
	uuid string

	ch    chan []byte
	conns chan *Conn

	mutex     sync.RWMutex
	listeners map[int]*Conn

	stats *metrics.Metrics

	id   int
	send int
}

var upgrader = websocket.Upgrader{} // use default options

func NewRelay(ws uint16, uuid string, stats *metrics.Metrics) *Relay {
	assert.NotNil(stats, "a metrics object must be provided")

	return &Relay{
		port: ws,
		uuid: uuid,

		ch:    make(chan []byte, 10),
		conns: make(chan *Conn, 10),

		mutex:     sync.RWMutex{},
		listeners: make(map[int]*Conn),

		stats: stats,

		id:   0,
		send: runtime.NumCPU(),
	}
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

	http.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		relay.stats.WritePrometheus(w)
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

func (relay *Relay) relayRange(listeners []*Conn, data []byte, wait *sync.WaitGroup) {
    for _, conn := range listeners {
		conn.msg(data)
	}

    wait.Done()
}

func (relay *Relay) relay(data []byte) {
	// quick write to prevent blocking if there is no listener
	select {
	case relay.ch <- data:
	default:
	}
	relay.mutex.RLock()
    msgCount := (len(relay.listeners) / relay.send + 1)
    wait := sync.WaitGroup{}

    curr := make([]*Conn, 0, msgCount)
    for _, conn := range relay.listeners {
        if len(curr) == msgCount {
            wait.Add(1)

            go relay.relayRange(curr, data, &wait)
            curr = make([]*Conn, 0, msgCount)
        }
        curr = append(curr, conn)
	}

    if len(curr) > 0 {
        wait.Add(1)
        go relay.relayRange(curr, data, &wait)
    }
    wait.Wait()
	relay.mutex.RUnlock()
}

func (relay *Relay) remove(id int) {
	relay.mutex.Lock()
	delete(relay.listeners, id)
	relay.mutex.Unlock()
	slog.Warn("removing client connection", "id", id)
	relay.stats.Set(metrics.CurrentConnections, len(relay.listeners))
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
    relay.mutex.Unlock()

	select {
	case relay.conns <- conn:
	default:
	}

	go conn.read()
	go conn.write()

	relay.stats.Set(metrics.CurrentConnections, len(relay.listeners))
	relay.stats.SetIfGreater(metrics.MaxConcurrentConnections, len(relay.listeners))
	relay.stats.Inc(metrics.TotalConnections)
}

func (relay *Relay) render(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// force sync on ids
	// probably should look into atomic ints here...
	relay.mutex.Lock()
	relay.id++
	id := relay.id
	relay.mutex.Unlock()

	relay.add(id, c)
	slog.Warn("connection established", "id", id, "addr", c.RemoteAddr())
}

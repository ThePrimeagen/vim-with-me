package main


import (
    "fmt"
    "net/http"
    "time"
)

// 1. Server side events (https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)
// 2. CSS grid stuff
// 3. Format going down to front end


func main() {
    http.HandleFunc("/events", eventsHandler)
    http.ListenAndServe(":8080", nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Simulate sending events (you can replace this with real data)
    for i := 0; i < 10; i++ {
        fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("Event %d", i))
        time.Sleep(2 * time.Second)
        w.(http.Flusher).Flush()
    }

    // Simulate closing the connection
    closeNotify := w.(http.CloseNotifier).CloseNotify()
    <-closeNotify
}

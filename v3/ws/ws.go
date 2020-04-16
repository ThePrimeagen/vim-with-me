package ws

import (
    "vim-with-me/types"
    "net/url"
    "net/http"
    "io"
    "io/ioutil"
    "log"
    "bytes"
    "fmt"
    "encoding/json"
    "github.com/gorilla/websocket"
)

type QuirkTokenRequest struct {
    Token string `json:"token"`
}

func createConnection(config map[string]string) *websocket.Conn {
    base := "https://websocket.quirk.gg"
    baseWs := "websocket.quirk.gg"

    var buf io.ReadWriter
    buf = new(bytes.Buffer)
    err := json.NewEncoder(buf).Encode(map[string]string{
        "auth_token": config["quirktoken"],
    })

    if err != nil {
        log.Fatalf("There is an error, and we cannot connect %+v\n", err)
    }

    request, err := http.NewRequest("POST", fmt.Sprintf("%s/token", base), buf);
    if (err != nil) {
        log.Fatalf("Couldn't make a new request: %+v\n", err)
    }
    request.Header.Add("Content-Type", "application/json");

    resp, err := http.DefaultClient.Do(request)

    if (err != nil) {
        log.Fatalf("Could not send request - %+v\n", err)
    }

    defer resp.Body.Close()

    var respData QuirkTokenRequest

    out, err := ioutil.ReadAll(resp.Body)

    if (err != nil) {
        log.Fatalf("Tried to read all %+v\n", err)
    }

    err = json.Unmarshal(out, &respData)

    if (err != nil) {
        log.Fatalf("THIS DIDNT WORK I BETS %d %+v\n", resp.StatusCode, err)
    }

    wsQuery := fmt.Sprintf("token=%s", respData.Token)

    u := url.URL{
        Scheme: "wss",
        Host: baseWs,
        Path: "",
        RawQuery: wsQuery,
    }

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatalln("dial:", err)
	}

    return c
}

type WS struct {
    Messages chan types.QuirkMessage
    Close chan interface{}
    config map[string]string
}

func CreateQuirk(config map[string]string) WS {
    c := createConnection(config)
    ws := WS{
        make(chan types.QuirkMessage),
        make(chan interface{}),
        config,
    }

    closed := false

    go func() {
        <- ws.Close
        closed = true
    }()

	go func() {
        defer c.Close()
		for {
            _, message, err := c.ReadMessage()

            if closed {
                return
            }

			if err != nil {
				log.Fatal("read:", err)
			}

            fmt.Printf("%s\n", message)
            var quirkMessage types.QuirkMessage
            err = json.Unmarshal(message, &quirkMessage)

            if (err != nil) {
                log.Printf("Unable to parse the quirk message. %s\n", message)
            }

            ws.Messages <- quirkMessage
		}
	}()

    log.Println("about to send over that irc info");
    log.Println("sent that info");

    return ws;
}


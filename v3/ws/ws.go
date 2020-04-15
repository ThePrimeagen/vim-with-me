package ws

import (
    "time"
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

type QuirkMessage struct {
    Data struct {
        Timestamp  time.Time `json:"timestamp"`
        Redemption struct {
            ID   string `json:"id"`
            User struct {
                ID          string `json:"id"`
                Login       string `json:"login"`
                DisplayName string `json:"display_name"`
            } `json:"user"`
            ChannelID  string    `json:"channel_id"`
            RedeemedAt time.Time `json:"redeemed_at"`
            Reward     struct {
                ID                  string      `json:"id"`
                ChannelID           string      `json:"channel_id"`
                Title               string      `json:"title"`
                Prompt              string      `json:"prompt"`
                Cost                int         `json:"cost"`
                IsUserInputRequired bool        `json:"is_user_input_required"`
                IsSubOnly           bool        `json:"is_sub_only"`
                Image               interface{} `json:"image"`
                DefaultImage        struct {
                    URL1X string `json:"url_1x"`
                    URL2X string `json:"url_2x"`
                    URL4X string `json:"url_4x"`
                } `json:"default_image"`
                BackgroundColor string `json:"background_color"`
                IsEnabled       bool   `json:"is_enabled"`
                IsPaused        bool   `json:"is_paused"`
                IsInStock       bool   `json:"is_in_stock"`
                MaxPerStream    struct {
                    IsEnabled    bool `json:"is_enabled"`
                    MaxPerStream int  `json:"max_per_stream"`
                } `json:"max_per_stream"`
                ShouldRedemptionsSkipRequestQueue bool        `json:"should_redemptions_skip_request_queue"`
                TemplateID                        interface{} `json:"template_id"`
            } `json:"reward"`
            Status string `json:"status"`
        } `json:"redemption"`
    } `json:"data"`
    Type string `json:"type"`
}

type WS struct {
    Messages chan QuirkMessage
    Close chan interface{}
    config map[string]string
}

func CreateQuirk(config map[string]string) WS {
    c := createConnection(config)
    ws := WS{
        make(chan QuirkMessage),
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
            var quirkMessage QuirkMessage
            err = json.Unmarshal(message, &quirkMessage)

            if (err != nil) {
                log.Printf("Unable to parse the quirk message. %s\n", message)
            }

            ws.Messages <- quirkMessage
		}
	}()

    go func() {
        ticker := time.NewTicker(30000 * time.Millisecond)
        defer ticker.Stop()

        for {
            <- ticker.C
            if closed {
                break
            }

            log.Printf("Sending that ping baby\n");
            c.WriteControl(
                websocket.PingMessage,
                []byte{},
                time.Now().Add(time.Second * 10))
        }
    }()

    log.Println("about to send over that irc info");
    log.Println("sent that info");

    return ws;
}


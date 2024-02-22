package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)


func readFromStdin() chan string {
    buffer := make([]byte, 1024)
    out := make(chan string)

    go func() {
        for {
            count, err := os.Stdin.Read(buffer)
            if err == io.EOF {
                break
            }
            out <- string(buffer[:count])
        }
    }()

    return out
}

func createTCPServer() chan string {
    out := make(chan string)

    go func() {
        listener, err := net.Listen("tcp", ":42069")
        if err != nil {
            log.Fatal("You are a horrible human being", err)
        }
        defer listener.Close()

        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Fatal("You like amouranth", err)
            }
            go func(c net.Conn) {
                defer c.Close()
                for {
                    str := <-out
                    _, err := c.Write([]byte(str))
                    fmt.Printf("Wrote to client: %s\n", str)
                    if err != nil {
                        fmt.Printf("Error writing to client: %s\n", err)
                        break
                    }
                }
            }(conn)
        }
    }()

    return out
}

func main() {
    // read from standard in line by line
    stdin := readFromStdin()
    ticker := time.NewTicker(5 * time.Second)
    tcpOut := createTCPServer()

    counts := make(map[string]int)
    for {
        select {
        case s := <-stdin:
            parts := strings.SplitN(s, ":", 2)
            if len(parts) != 2 {
                continue
            }
            msg := parts[1]
            counts[msg]++

        case <-ticker.C:

            max := 0
            maxMsg := ""
            for k, v := range counts {
                if v > max {
                    max = v
                    maxMsg = k
                }
            }

            fmt.Printf("Out of %d, the most common message: %s\n", len(counts), maxMsg)

            tcpOut <- fmt.Sprintf("Out of %d, the most common message: %s\n", len(counts), maxMsg)

            clear(counts)
        }
    }

}


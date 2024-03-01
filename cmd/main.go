package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"chat.theprimeagen.com/pkg/processors"
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

//var allowableChars = []string{"<dot>", "<backspace>", "<space>", "<esc>", "<cr>", "<tab>"};
//func contains(arr []string, str string) bool {
//    for _, a := range arr {
//        if a == str {
//            return true
//        }
//    }
//    return false
//}
//
func main() {
    // read from standard in line by line
    stdin := readFromStdin()
    //tcpOut := createTCPServer()

    processor := processors.NewTDProcessor(5)

    for {
        select {
        case s := <-stdin:
            processor.Process(strings.TrimSpace(s))

        case point := <-processor.Out():
            fmt.Printf("Got a point: %s\n", point)
            //tcpOut <- point
        }
    }

}


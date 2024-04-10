package main

import (
	"fmt"
	"log"
	"net"
)


func main() {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 42073))
    if err != nil {
        log.Fatalf("you suck %+v\n", err)
    }

    for {
		conn, err := listener.Accept()
        if err != nil {
            log.Fatalf("here is that server error: %+v\n", err)
        }


    }
}

package main

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

func read_conn(conn Connection) {
	for {
		log.Printf("Reading(%d)...\n", id)
        cmd, err := readTCPCommand(&conn.FrameReader)
        if err != nil {
            log.Printf("error with: %+v\n", err)
        }

		log.Printf("got command %+v", cmd)
        err = conn.Write(cmd)

        if err != nil {
            log.Printf("error sending command back: %+v\n", err)
			break
        }
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 42073))
	defer listener.Close()
	if err != nil {
		log.Fatalf("you suck %+v\n", err)
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("starting server\n")
	fmt.Printf("starting server from fmt!\n")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalf("here is that server error: %+v\n", err)
		}

        newConn := NewConnection(conn)
		go read_conn(conn, myId)
		go write_conn(conn, myId)
	}
}

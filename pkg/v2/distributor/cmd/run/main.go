package main

import (
	"flag"
	"github.com/theprimeagen/vim-with-me/pkg/v2/distributor"
	"os"
	"strconv"
	"log/slog"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "port to listen on")
	flag.Parse()
	downstreams := os.Args[1:]

	if port == 0 {
		portStr := os.Getenv("PORT")
		if portStr == "" {
			slog.Error("No port specified!")
			os.Exit(1)
		}
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			slog.Error("Error converting port to int", "port", portStr, "err", err)
			os.Exit(1)
		}
	}

	authId := os.Getenv("AUTH_ID")
	if authId == "" {
		slog.Error("No auth id specified!")
		os.Exit(1)
	}

	d := distributor.NewDistributor(port, authId, downstreams)
	d.Run()
}

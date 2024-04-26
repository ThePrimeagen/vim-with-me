package chat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type ChatMsg struct {
	Name string
	Msg  string
	Bits int
}

func parse(msg string) (*ChatMsg, error) {
	parts := strings.SplitN(msg, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("malformed message")
	}

	messageType := parts[0]
	switch messageType {
	case "message":
		msgParts := strings.SplitN(parts[1], ":", 2)
		return &ChatMsg{
			Name: msgParts[0],
			Msg:  msgParts[1],
			Bits: 0,
		}, nil

	case "bits":
		msgParts := strings.SplitN(parts[1], ":", 3)

		bits, err := strconv.Atoi(msgParts[1])

		if err != nil {
			return nil, err
		}

		return &ChatMsg{
			Name: msgParts[0],
			Bits: bits,
			Msg:  msgParts[2],
		}, nil
	}

	return nil, errors.New("unknown message type")
}

type Chat struct {
	Messages chan ChatMsg
	path     string
	stdout   io.ReadCloser
    done     chan struct{}
}

func NewChat(path string) Chat {
	return Chat{
		path:     path,
		Messages: make(chan ChatMsg),
        done:     make(chan struct{}),
	}
}

func (c *Chat) readFromStdout() {
	buf_reader := bufio.NewReader(c.stdout)
	for {
		line, _, err := buf_reader.ReadLine()
		if err != nil {
			// TODO: how to close the thing properly
			c.done <- struct{}{}
			close(c.Messages)
			break
		}

		msg, err := parse(string(line))
		if err != nil {
			continue
		}

		c.Messages <- *msg
	}

}

func (c *Chat) Start() (chan struct{}, error) {
	// spawn program and read from the standard out

	cmd := exec.Command(c.path)
	defer cmd.Process.Kill()

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

    c.stdout = stdout

	if err := cmd.Start(); err != nil {
		return nil, err
	}

    go c.readFromStdout()

	return c.done, nil
}

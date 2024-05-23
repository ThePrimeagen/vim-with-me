package controller

import (
	"context"
	"time"
)

type controllerChan chan time.Time
type Next interface {
	Next() string
}
type SendKey interface {
	SendKey(string) string
}

type Controller struct {
	next      Next
	send      SendKey
    curr      string
	nextInput controllerChan
	playInput controllerChan
}

func NewController(next Next, send SendKey) *Controller {
	return &Controller{
        next: next,
        send: send,
        curr: "",
    }
}

func (c *Controller) WithNextInput(nextInput controllerChan) *Controller {
	c.nextInput = nextInput
	return c
}

func (c *Controller) WithPlayInput(playInput controllerChan) *Controller {
	c.playInput = playInput
	return c
}

func (c *Controller) Start(ctx context.Context) {
outer:
	for {
		select {
		case <-c.playInput:
            if c.curr == "" {
                continue
            }
            c.send.SendKey(c.curr)
		case <-c.nextInput:
            c.curr = c.next.Next()
		case <-ctx.Done():
			break outer
		}

	}
}

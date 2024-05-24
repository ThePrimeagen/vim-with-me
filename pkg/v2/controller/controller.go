package controller

import (
	"context"
	"fmt"
	"time"
)

type controllerChan <-chan time.Time
type Next interface {
	Next() string
}
type SendKey interface {
	SendKey(string)
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

func (c *Controller) WithInputTimer(nextInput controllerChan) *Controller {
	c.nextInput = nextInput
	return c
}

func (c *Controller) WithPlayTimer(playInput controllerChan) *Controller {
	c.playInput = playInput
	return c
}

func (c *Controller) Start(ctx context.Context) {
outer:
	for {
		select {
		case <-c.playInput:
            if c.curr == "" {
                break
            }
            c.send.SendKey(c.curr)
		case <-c.nextInput:
            c.curr = c.next.Next()
		case <-ctx.Done():
			break outer
		}

	}
}

func (c *Controller) String() string {
    return fmt.Sprintf("Controller: %s", c.curr)
}

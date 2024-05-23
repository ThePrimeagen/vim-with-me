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
	nextInput controllerChan
	playInput controllerChan
}

func NewController(next Next, send SendKey) *Controller {
	return &Controller{
        next: next,
        send: send,
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
		case <-c.nextInput:
		case <-ctx.Done():
			break outer
		}

	}
}

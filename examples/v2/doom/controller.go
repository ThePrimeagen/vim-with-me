package doom

import (
	"context"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
)

var timeBetweenUse = time.Millisecond * 500
type DoomController struct {
	send       controller.SendKey
	controller *controller.Controller

    timeSinceLastUse time.Time
}

func NewDoomController(next controller.Next, send controller.SendKey) {
    dc := DoomController{}
    dc.controller = controller.NewController(next, &dc)
    dc.timeSinceLastUse = time.Now()
}

func (dc *DoomController) Start(ctx context.Context) {
    go dc.controller.Start(ctx)
}

func (dc *DoomController) SendKey(key string) {
    now := time.Now()
    for _, k := range key {
        if k == 'e' {
            if dc.timeSinceLastUse.Sub(now) >= timeBetweenUse {
                dc.timeSinceLastUse = now
            } else {
                continue
            }
        }
        dc.send.SendKey(string(k))
    }
}


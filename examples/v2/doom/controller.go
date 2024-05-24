package doom

import (
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
)

type DoomController struct {
	send controller.SendKey

	timeSinceLastUse time.Time
	timeBetweenUse   time.Duration
}

func NewDoomController(send controller.SendKey) *DoomController {
	dc := DoomController{}
	dc.timeSinceLastUse = time.Now()
	dc.timeBetweenUse = 500
	dc.send = send
    return &dc
}

func (dc *DoomController) WithTimeBetweenUse(useTime time.Duration) *DoomController {
	dc.timeBetweenUse = useTime
	return dc
}

func (dc *DoomController) Init() *DoomController {
	return dc
}

func (dc *DoomController) SendKey(key string) {
	now := time.Now()
	for _, k := range key {
		if k == 'e' {
			if now.Sub(dc.timeSinceLastUse) >= dc.timeBetweenUse {
				dc.timeSinceLastUse = now
			} else {
				continue
			}
		}
		dc.send.SendKey(string(k))
	}
}

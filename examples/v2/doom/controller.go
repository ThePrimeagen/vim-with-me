package doom

import (
	"slices"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
)

type DoomController struct {
	send controller.SendKey

    playing bool

	timeSinceLastUse time.Time
	timeBetweenUse   time.Duration
}

func NewDoomController(send controller.SendKey) *DoomController {
    return &DoomController{
        timeSinceLastUse: time.Now(),
        timeBetweenUse: 500,
        send: send,
        playing: false,
    }
}

func (dc *DoomController) Play() {
    dc.playing = true
}

func (dc *DoomController) Stop() {
    dc.playing = false
}

func (dc *DoomController) WithTimeBetweenUse(useTime time.Duration) *DoomController {
	dc.timeBetweenUse = useTime
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

        if dc.playing {
            dc.send.SendKey(string(k))
        }
	}
}

func DoomChatMapFn(msg string) string {
	if len(msg) > 2 {
		return msg
	}

    msg = strings.ToLower(msg)
	bytes := []byte(msg)
	slices.Sort(bytes)

	return string(bytes)
}

var validMsg = "wasdfeasdwfwawaedfdeeffsesewdsaf"

func DoomFilterFn(msg string) bool {
	if len(msg) > 2 {
		return false
	}
	return strings.Contains(validMsg, msg)
}

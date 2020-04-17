package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"vim-with-me/config"
	"vim-with-me/nvim"
	"vim-with-me/system"
	"vim-with-me/ws"
)

var lastTimeUsed map[string]int64
var THROTTLE_TIME int64

func setThrottle(userId string) bool {
	now := time.Now()
	unixNano := now.UnixNano()
	milli := unixNano / 1000000

	if val, ok := lastTimeUsed[userId]; ok {
		if milli-val < THROTTLE_TIME {
			return true
		}
	}

	lastTimeUsed[userId] = milli
	return false
}

func execSystemCommand(cmd system.SystemCommand) {
	fmt.Printf("About to execute %s %+v\n", cmd.Program, cmd.OnArgs)
	sysCmd := exec.Command(cmd.Program, cmd.OnArgs...)

	sysCmd.Stdin = os.Stdin
	sysCmd.Stdout = os.Stdout
	sysCmd.Stderr = os.Stderr
	err := sysCmd.Run()

	if err == nil {
		fmt.Printf("execSystemCommand#err %+v\n", err)
	}

	<-cmd.Scheduled.Done

	fmt.Printf("About to execute %s %+v\n", cmd.Program, cmd.OffArgs)
	sysCmd = exec.Command(cmd.Program, cmd.OffArgs...)
	sysCmd.Run()
}

func main() {
	THROTTLE_TIME = 1000

	asdf := system.CreateASDF()
	xrandr := system.CreateXrandr()

	c := config.ReadConfig("./.config")
	quirk := ws.CreateQuirk(c)
	n := nvim.CreateVimWithMe()
	lastTimeUsed = make(map[string]int64)

	for {
		msg := <-quirk.Messages
		name := msg.Data.Redemption.Reward.Title
		userId := msg.Data.Redemption.User.ID

		if setThrottle(userId) {
			fmt.Printf("You have been throttled %s\n", msg.Data.Redemption.User.DisplayName)
			continue
		}

		switch name {
		case "Vim Command":

			// Do we want to parse it better?
			n.SendCommand <- nvim.NvimCommand{
				msg.Data.Redemption.User.ID,
				msg,
			}

		case "Xrandr":
			xrandr.Scheduled.AddSeconds(5)
			go execSystemCommand(xrandr)

		case "asdf":
			asdf.Scheduled.AddSeconds(3)
			go execSystemCommand(asdf)
		case "Vim ColorScheme":
			n.SendCommand <- nvim.NvimColor{
				msg.Data.Redemption.User.ID,
				msg,
			}
		default:
			fmt.Printf("Received %s -- Unable to process\n", name)
		}

		n.SendCommand <- msg
	}
}

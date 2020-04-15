package nvim

import (
    "log"
    "os"
    "fmt"
    "github.com/neovim/go-client/nvim"
)


func connect(addr string) *nvim.Nvim {

	// Dial with default options.
	v, err := nvim.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

    return v
}

type NvimCommand struct {
    UserId string
    Type string
    Command string
    From string
}

type NvimExecutor struct {
    SendCommand chan interface{}
    nvim *nvim.Nvim
}

func executeCommand(command NvimCommand) {

    // Todo this.
}

func CreateVimWithMe() NvimExecutor {
    addr := os.Getenv("NVIM_LISTEN_ADDRESS")
	if addr == "" {
		log.Fatal("NVIM_LISTEN_ADDRESS not set")
	}

    v := connect(addr)

    nvimExec := NvimExecutor{
        make(chan interface{}),
        v,
    }

	go func() {
        for {
            command := <- nvimExec.SendCommand

            switch t := command.(type) {
            case string:
                if command == "close" {
                    defer v.Close()
                    return;
                }

            case NvimCommand:
                //executeCommand(t)

            default:
                fmt.Printf("I don't know about type %T!\n", t)
            }
        }
    }()

    return nvimExec
}


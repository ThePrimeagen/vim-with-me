package nvim

import (
    "log"
    "os"
    "fmt"
    "vim-with-me/types"
    "github.com/neovim/go-client/nvim"
)

var allowedVimCommands []string

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
    Input types.QuirkMessage
}


type NvimColor struct {
    UserId string
    Input types.QuirkMessage
}

type NvimExecutor struct {
    SendCommand chan interface{}
    nvim *nvim.Nvim
}
/*
{
  "data": {
    "timestamp": "2020-04-16T01:06:04.897350104Z",
    "redemption": {
      "id": "b6a561b5-176b-4102-afae-6775cabd20eb",
      "user": {
        "id": "167160215",
        "login": "theprimeagen",
        "display_name": "ThePrimeagen"
      },
      "channel_id": "167160215",
      "redeemed_at": "2020-04-16T01:06:04.812304698Z",
      "reward": {
        "id": "cc97679e-b8db-433b-86e5-6a9c1e2c76db",
        "channel_id": "167160215",
        "title": "Vim Command",
        "prompt": "Send a vim command",
        "cost": 69,
        "is_user_input_required": true,
        "is_sub_only": false,
        "image": null,
        "default_image": {
          "url_1x": "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png",
          "url_2x": "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png",
          "url_4x": "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"
        },
        "background_color": "#BDA8FF",
        "is_enabled": true,
        "is_paused": false,
        "is_in_stock": true,
        "max_per_stream": {
          "is_enabled": false,
          "max_per_stream": 0
        },
        "should_redemptions_skip_request_queue": true,
        "template_id": null
      },
      "user_input": "20j",
      "status": "FULFILLED"
    }
  },
  "type": "TWITCH_CHANNEL_REWARD"
}
*/
func getRepeatAmount(input string) int {
    start := -1

    for i, c := range input {
        if c < '0' || c > '9' {
            break
        }
        start = i
	}

    return start + 1
}

func executeColor(vim *nvim.Nvim, command NvimColor) {
    input := command.Input.Data.Redemption.UserInput

    // I feel like someone would do something naughty
    if len(input) > 20 {
        fmt.Printf("Error#executeCommand input was greater than 20 :: %s\n", input)
        return
    }

    fmt.Printf("Executing vim color %s\n", input)
    vim.Command(fmt.Sprintf("colorscheme %s", input))
}

func executeCommand(vim *nvim.Nvim, command NvimCommand) {
    input := command.Input.Data.Redemption.UserInput
    start := getRepeatAmount(input)

    fmt.Printf("executeCommand: %s %d\n", input, start)

    if !contains(allowedVimCommands, input[start:]) {
        fmt.Printf("Error#executeCommand Did not find %s in %+v\n", input[start:], allowedVimCommands)
        return
    }

    // I feel like someone would do something naughty
    if len(input) > 6 {
        fmt.Printf("Error#executeCommand input was greater than 6 :: %s\n", input[start:])
        return
    }

    fmt.Printf("Executing vim command %s\n", input)
    vim.Command(fmt.Sprintf("normal!%s", input))
}

func contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
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
                executeCommand(v, t)

            case NvimColor:
                executeColor(v, t)

            default:
                fmt.Printf("I don't know about type %T!\n", t)
            }
        }
    }()

    return nvimExec
}

func init() {
    allowedVimCommands = []string{
        "j",
        "k",
        "h",
        "l",
        "gg",
        "G",
        "H",
        "L",
        "zz",
        "V",
        "~",
        "%",
        "I",
        "A",
        "v",
        "zt",
        "zb",
        "gv",
        "gi",
        "Vjj",
    }
}


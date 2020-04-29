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

func getCommandStartIdx(input string) int {
    start := -1

    for i, c := range input {
        if c < '0' || c > '9' {
            break
        }
        start = i
	}

    return start + 1
}

func isValid(s string) bool {
    A := int('A')
    a := int('a')
    Z := int('Z')
    z := int('z')
    dunder := int('_')
    dash := int('-')

    valid := true
    for _, c := range s {
        char := int(c);
        valid = char >= 48 && char <= 57 ||
            char >= a && char <= z || char >= A && char <= Z ||
            char == dunder || char == dash

        if !valid {
            break;
        }
    }

    return valid
}

func executeColor(vim *nvim.Nvim, command NvimColor) {
    input := command.Input.Data.Redemption.UserInput

    if !isValid(input) {
        fmt.Printf("Not valid input ya dingus %s \n", input)
        return
    }

    fmt.Printf("Executing vim color %s\n", input)
    vim.Command(fmt.Sprintf("colorscheme %s", input))
}

func executeCommand(vim *nvim.Nvim, command NvimCommand) {
    input := command.Input.Data.Redemption.UserInput
    start := getCommandStartIdx(input)

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
        "vjj",
    }
}


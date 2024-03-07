package commands

import "chat.theprimeagen.com/pkg/tcp"

func Render(data string) tcp.TCPCommand {
    return tcp.TCPCommand{
        Command: "r",
        Data: data,
    }
}

func Close(msg string) tcp.TCPCommand {
    return tcp.TCPCommand{
        Command: "c",
        Data: msg,
    }
}

func Error(msg string) tcp.TCPCommand {
    return tcp.TCPCommand{
        Command: "e",
        Data: msg,
    }
}


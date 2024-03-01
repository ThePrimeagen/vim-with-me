package processors

type Processor interface {
    Process(str string)
    Out() chan string
}


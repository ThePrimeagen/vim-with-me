package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
	"github.com/theprimeagen/vim-with-me/pkg/v2/cmd"
)

func getDebug(name string) (*testies.DebugFile, error) {
    if name == "" {
        return testies.EmptyDebugFile(), nil
    }
    return testies.NewDebugFile(name)
}

type CmdErrParser struct {
    debug *testies.DebugFile
    round chan int
    prompt chan string
    done chan string
    doneStr string
    doneReadCount int
    promptStr string
    readingPrompt bool
    empty bool
}

func newCmdErrParser(debug *testies.DebugFile) CmdErrParser {
    return CmdErrParser{
        debug: debug,
        round: make(chan int, 1),
        prompt: make(chan string, 1),
        done: make(chan string, 1),
        doneReadCount: 0,
        promptStr: "",
        doneStr: "",
        readingPrompt: false,
        empty: false,

    }
}

func (c *CmdErrParser) parseWait(countStr string) {
    countStrTrimmed := strings.TrimSpace(countStr)
    count, err := strconv.Atoi(countStrTrimmed)
    assert.NoError(err, "zig program gave me a non number", "parts", countStr)

    c.round <- count
}

func (c *CmdErrParser) parse(b []byte) (int, error) {
    str := string(b)

    lines := strings.Split(str, "\n")
    idx := 0

    for idx < len(lines) {
        line := lines[idx]

        if c.readingPrompt {
            for idx < len(lines) {
                curr := lines[idx]
                idx++

                isEmpty := curr == ""
                if isEmpty && c.empty {
                    c.prompt <- c.promptStr
                    c.readingPrompt = false
                    break;
                }

                c.empty = isEmpty
                if c.empty {
                    continue
                }

                c.promptStr += curr + "\n"
            }

        } else if strings.Contains(line, "waiting") {
            parts := strings.Split(line, "-")
            assert.Assert(len(parts) == 2, "somehow i didn't communicate back with 2 pieces", "parts", parts)
            c.parseWait(parts[1])
        } else if line == "prompt:" {
            c.readingPrompt = true
            c.promptStr = ""
            c.empty = false
        } else if line == "final" {
            for range 2 {
                idx++
                if idx >= len(lines) {
                    break;
                }
                c.doneReadCount++
                c.doneStr += lines[idx] + "\n"
            }

            if c.doneReadCount == 2 {
                c.done <- c.doneStr
            }
        } else {
            c.debug.WriteLine([]byte(fmt.Sprintf("td: %s\n", line)))
        }

        idx++
    }

    return len(b), nil
}

type OpenAIChat struct {
    client *openai.Client
    system string
    ctx context.Context
}

func (o *OpenAIChat) chat(chat string) (string, error) {
    resp, err := o.client.CreateChatCompletion(o.ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: o.system,
            },
            {
                Role: openai.ChatMessageRoleUser,
                Content: chat,
            },
        },
    })
    if err != nil {
        return "", nil
    }

    return resp.Choices[0].Message.Content, nil
}

func NewOpenAIChat(secret string, system string, ctx context.Context) *OpenAIChat {
    client := openai.NewClient(secret)
    return &OpenAIChat{
        system: system,
        client: client,
        ctx: ctx,
    }
}

type StatefulOpenAIChat struct {
    ai *OpenAIChat
    out chan string
    ctx context.Context
}

func newStatefulOpenAIChat(secret string, system string, ctx context.Context) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(secret, system, ctx),
        out: make(chan string, 10),
        ctx: ctx,
    }
}

func (s *StatefulOpenAIChat) prompt(p string) error {
    str, err := s.ai.chat(p)
    if err != nil {
        return err
    }

    s.out <- str
    return nil
}

func (s *StatefulOpenAIChat) ReadWithTimeout(timeout time.Duration) string {
    select {
    case <-s.ctx.Done():
        return ""
    case <-time.NewTimer(timeout).C:
        return ""
    case resp := <-s.out:
        return resp
    }
}

type RandomPos struct {
    outOfBounds string
}

func NewRandomPos(maxRows int) RandomPos {
    return RandomPos{outOfBounds: fmt.Sprintf("%d,69", maxRows + 1)}
}

// As of right now, out of bounds guesses will place a tower within your area
// randomly
func (r *RandomPos) NextPos() string {
    return r.outOfBounds
}

type GameState struct {
}

type PosGenerator interface {
    GetPositions(count int, promptState GameState) []string
}

// if the ai provides some line that is unreadable i will just provide a value too large and thus random placement
func getPosFromAIResponse(line string) string {
    parts := strings.Split(line, ",")
    if len(parts) != 2 {
        return "999,999"
    }

    _, err := strconv.Atoi(parts[0])
    if err != nil {
        return "999,999"
    }

    _, err = strconv.Atoi(parts[1])
    if err != nil {
        return "999,999"
    }
    return line
}

func main() {

    testies.SetupLogger()

	godotenv.Load()

	debugFile := ""
	flag.StringVar(&debugFile, "debug", "", "runs the file like the program instead of running doom")

	systemPromptFile := "PROMPT"
	flag.StringVar(&systemPromptFile, "system", "PROMPT", "the system prompt to use")

	seed := 1337
	flag.IntVar(&seed, "seed", 69420, "the seed value for the game")
	flag.Parse()


	args := flag.Args()
    assert.Assert(len(args) >= 2, "you must provide path to exec and json file")
	name := args[0]
	json := args[1]

    debug, err := getDebug(debugFile)
    if err != nil {
        log.Fatalf("could not open up debug file: %v\n", err)
    }
    defer debug.Close()

    systemPrompt, err := os.ReadFile(systemPromptFile)
    if err != nil {
        log.Fatalf("could not open system prompt: %+v\n", err)
    }

	ctx := context.Background()
    ctxWithCancel, cancel := context.WithCancel(ctx)

	twitchChat, err := chat.NewTwitchChat(ctx)
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(td.TDFilter(24, 80));

    errParser := newCmdErrParser(debug)
    cmdr := cmd.NewCmder(name, ctxWithCancel).
        AddVArg(json).
        AddKVArg("--seed", fmt.Sprintf("%d", seed)).
        WithErrFn(errParser.parse).
        WithOutFn(func(b []byte) (int, error) {
            //fmt.Printf("%s\n", string(b))
            return len(b), nil
        })

    go cmdr.Run()
	go chtAgg.Pipe(twitchChat)

    ai := newStatefulOpenAIChat(os.Getenv("OPENAI_API_KEY"), string(systemPrompt), ctx)

    done := make(chan struct{}, 1)
    go func() {
        for doneStr := range errParser.done {
            debug.WriteLine([]byte("------------------------------------------"))
            debug.WriteLine([]byte("-------------- game over -----------------"))
            debug.WriteLine([]byte("------------------------------------------"))
            debug.WriteLine([]byte(doneStr))
            cancel();

            // note: waiting a moment to prevent any more IO to screw up my printing
            <-time.NewTimer(time.Millisecond * 250).C
            debug.WriteStrLine(fmt.Sprintf("Results: %s", doneStr))

            if strings.Contains(doneStr, "1: won") {
                fmt.Printf("%s %d 1\n", systemPromptFile, seed)
            } else {
                fmt.Printf("%s %d 2\n", systemPromptFile, seed)
            }

            done<-struct{}{}
        }
    }()

    go func() {
        outer:
        for {
            select {
            case p := <-errParser.prompt:
                debug.WriteLine([]byte("prompt received"))
                err := ai.prompt(p)
                assert.NoError(err, "error when prompting", "err", err)
            case <-ctxWithCancel.Done():
                break outer;
            }
        }
    }()

    defer func() {
        fmt.Println("\x1b[?25h")
    }()

    posGen := NewRandomPos(24)
    roundCount := 0
    outer:
    for {
        debug.WriteStrLine(fmt.Sprintf("------------- waiting on game round: %d -----------------", roundCount))
        select {
        case <-ctxWithCancel.Done():
            break outer;
        case count := <- errParser.round:
            resp := ai.ReadWithTimeout(time.Second * 15)
            assert.Assert(resp != "", "i got an empty response from the ai???")
            roundCount++

            sent := 0
            for _, line := range strings.Split(resp, "\n") {
                if line == "" || sent == count {
                    break;
                }

                line = getPosFromAIResponse(line)

                sent++
                debug.WriteStrLine(fmt.Sprintf("ai-placement(%d): %s\n", roundCount, line))
                err := cmdr.WriteLine([]byte(fmt.Sprintf("2%s\n", line)))
                if err != nil {
                    debug.WriteStrLine(fmt.Sprintf("error writing ai coord: %+v", err))
                }
                assert.NoError(err, "error trying to write line to program", "err", err)
            }

            assert.Assert(sent == count, "we sent incorrect amount of positions", "sent", sent, "count", count)

            for range count {
                cmdr.WriteLine([]byte(fmt.Sprintf("1%s\n", posGen.NextPos())))
            }
        }
    }


    <-done
}


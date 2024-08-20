package main

import (
	"context"
	"encoding/json"
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
    gs chan GameState
    done chan string
    readingPrompt bool
    empty bool
}

func newCmdErrParser(debug *testies.DebugFile) CmdErrParser {
    return CmdErrParser{
        debug: debug,
        gs: make(chan GameState, 1),
        done: make(chan string, 1),
        readingPrompt: false,
        empty: false,

    }
}

func (c *CmdErrParser) parse(b []byte) (int, error) {
    var gs GameState;
    err := json.Unmarshal(b, &gs)
    if err != nil {
        fmt.Printf("td: %s\n", string(b))
    }
    c.gs <- gs

    return len(b), nil
}

type OpenAIChat struct {
    client *openai.Client
    system string
}

var foo = 1337

func (o *OpenAIChat) chat(chat string, ctx context.Context) (string, error) {
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini20240718,
        Temperature: 0.55,
        Seed: &foo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: o.system,
            },
            {
                Role: openai.ChatMessageRoleSystem,
                Content: `Response Format MUST BE line separated Row,Col tuples
Example output with 2 positions specified
20,3
12,4

This specified position of row 20 col 3 and second, separated by new line, row 12 and col 4
`,
            },
            {
                Role: openai.ChatMessageRoleUser,
                Content: chat,
            },
        },
    })
    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}

func NewOpenAIChat(secret string, system string) *OpenAIChat {
    client := openai.NewClient(secret)
    return &OpenAIChat{
        system: system,
        client: client,
    }
}

type StatefulOpenAIChat struct {
    ai *OpenAIChat
    ctx context.Context
}

func newStatefulOpenAIChat(secret string, system string, ctx context.Context) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(secret, system),
        ctx: ctx,
    }
}

func (s *StatefulOpenAIChat) prompt(p string, ctx context.Context) (string, error) {
    str, err := s.ai.chat(p, ctx)
    if err != nil {
        return "", err
    }
    return str, err
}

func (s *StatefulOpenAIChat) ReadWithTimeout(p string, t time.Duration) (string, error) {
    ctx, cancel := context.WithCancel(s.ctx)
    go func() {
        <-time.NewTimer(t).C
        cancel()
    }()

    return s.ai.chat(p, ctx)
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

type BoxPos struct {
    maxRows int
    position int
}

func NewBoxPos(maxRows int) BoxPos {
    return BoxPos{
        maxRows: maxRows,
        position: 0,
    }
}

// As of right now, out of bounds guesses will place a tower within your area
// randomly
func (r *BoxPos) NextPos() string {
    col := 6
    if r.position & 0x1 == 0 {
        col = 12
    }

    row := ((r.position % 4) / 2) * 5
    r.position++

    return fmt.Sprintf("%d,%d", row, col)
}


type PositionGenerator interface  {
    NextPos() string
}

type Tower struct {
    Row int `json:"row"`
    Col int `json:"col"`
    Ammo int `json:"ammo"`
    Level int `json:"level"`
}

type Range struct {
    StartRow uint `json:"startRow"`
    EndRow uint `json:"endRow"`
}

type GameState struct {
    Rows uint `json:"rows"`
    Cols uint `json:"cols"`
    AllowedTowers int `json:"allowedTowers"`
    YourCreepDamage uint `json:"yourCreepDamage"`
    EnemyCreepDamage uint `json:"enemyCreepDamage"`
    YourTowers []Tower `json:"yourTowers"`
    EnemyTowers []Tower `json:"enemyTowers"`
    TowerPlacementRange Range `json:"towerPlacementRange"`
    CreepSpawnRange Range `json:"creepSpawnRange"`
    Round uint `json:"round"`
    Finished bool `json:"finished"`
    Playing bool `json:"playing"`
    Winner uint `json:"winner"`
}

func (gs *GameState) String() string {
    b, err := json.Marshal(gs)
    assert.NoError(err, "unable to create gamestate string")
    return string(b)
}

type PosGenerator interface {
    GetPositions(count int, promptState GameState) []string
}

// if the ai provides some line that is unreadable i will just provide a value too large and thus random placement
func getPosFromAIResponse(line string) string {
    parts := strings.Split(line, ",")
    for i, p := range parts {
        parts[i] = strings.TrimSpace(p)
    }

    if len(parts) != 2 {
        return "999,999"
    }

    row, err := strconv.Atoi(parts[0])
    if err != nil {
        return "999,999"
    }

    col, err := strconv.Atoi(parts[1])
    if err != nil {
        return "999,999"
    }
    return fmt.Sprintf("%d,%d", row, col)
}

func main() {

    testies.SetupLogger()

	godotenv.Load()

	debugFile := ""
	flag.StringVar(&debugFile, "debug", "", "runs the file like the program instead of running doom")

	systemPromptFile := "THEPRIMEAGEN"
	flag.StringVar(&systemPromptFile, "system", "THEPRIMEAGEN", "the system prompt to use")

	viz := false
	flag.BoolVar(&viz, "viz", false, "displays the game")


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
	twitchChat, err := chat.NewTwitchChat(ctx)
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(td.TDFilter(24, 80));

    errParser := newCmdErrParser(debug)
    cmdr := cmd.NewCmder(name, ctx).
        AddVArg(json).
        AddKVArg("--seed", fmt.Sprintf("%d", seed)).
        WithErrFn(errParser.parse).
        WithOutFn(func(b []byte) (int, error) {
            if viz {
                fmt.Printf("%s\n", string(b))
            }
            return len(b), nil
        })

    go cmdr.Run()
	go chtAgg.Pipe(twitchChat)

    ai := newStatefulOpenAIChat(os.Getenv("OPENAI_API_KEY"), string(systemPrompt), ctx)
    aiRandomGuesses := 0
    aiBadParse := 0
    aiTotalTowers := 0
    roundCount := 0
    fmt.Printf("won,round,prompt file,seed,ai total towers,ai guesses,ai bad parses\n")

    defer func() {
        fmt.Println("\x1b[?25h")
    }()

    posGen := NewBoxPos(24)
    outer:
    for {
        debug.WriteStrLine(fmt.Sprintf("------------- waiting on game round: %d -----------------", roundCount))
        select {
        case <-ctx.Done():
            break outer;
        case gs := <- errParser.gs:
            debug.WriteStrLine(fmt.Sprintf("ai-placement response: \"%s\"", gs.String()))
            count := gs.AllowedTowers
            aiResponses := 0
            tries := 0
            roundCount = int(gs.Round)

            if gs.Finished {
                if gs.Winner == '1' {
                    fmt.Printf("1,%d,%s,%d,%d,%d,%d\n", roundCount, systemPromptFile, seed, aiTotalTowers, aiRandomGuesses, aiBadParse)
                } else {
                    fmt.Printf("2,%d,%s,%d,%d,%d,%d\n", roundCount, systemPromptFile, seed, aiTotalTowers, aiRandomGuesses, aiBadParse)
                }
                break outer
            }

            if gs.Playing {
                continue
            }

            for aiResponses < count && tries < 3 {
                resp, err := ai.ReadWithTimeout(gs.String(), time.Second * 5)
                tries++

                if err != nil {
                    debug.WriteStrLine(fmt.Sprintf("ai-placement response: \"%s\" err: \"%s\"", resp, err.Error()))
                    parts := strings.Split(err.Error(), "try again in ")
                    if len(parts) == 2 {
                        secsStr := strings.Split(parts[1], " ")[0]
                        secs, err := strconv.ParseFloat(secsStr[0:len(secsStr) - 2], 64)
                        if err == nil {
                            dur := time.Duration(float64(time.Second) * secs)
                            debug.WriteStrLine(fmt.Sprintf("ai-placement wait time required: %d", dur))
                            <-time.NewTimer(dur).C
                        }
                    }
                    continue
                }

                if resp == "" {
                    <-time.NewTimer(time.Second).C
                }

                for _, line := range strings.Split(resp, "\n") {
                    line = strings.TrimSpace(line)
                    if line == "" || aiResponses == count {
                        break;
                    }

                    parsedLine := getPosFromAIResponse(line)
                    debug.WriteStrLine(fmt.Sprintf("ai-placement: %s - %s", line, parsedLine))
                    if parsedLine == "999,999" {
                        aiBadParse++
                        continue
                    }

                    aiTotalTowers++
                    aiResponses++

                    debug.WriteStrLine(fmt.Sprintf("ai-placement(%d): %s\n", roundCount, parsedLine))
                    err := cmdr.WriteLine([]byte(fmt.Sprintf("2%s\n", parsedLine)))
                    if err != nil {
                        debug.WriteStrLine(fmt.Sprintf("error writing ai coord: %+v", err))
                    }
                    assert.NoError(err, "error trying to write line to program", "err", err)
                }

            }

            for aiResponses < count {
                err := cmdr.WriteLine([]byte("2999,999"))
                if err != nil {
                    debug.WriteStrLine(fmt.Sprintf("error writing ai coord: %+v", err))
                }
                aiResponses++
                aiRandomGuesses++
                aiTotalTowers++
            }

            for range count {
                next := fmt.Sprintf("1%s\n", posGen.NextPos())
                debug.WriteStrLine(fmt.Sprintf("PLAYER MOVE: %s", next))
                cmdr.WriteLine([]byte(next))
            }
        }
    }
}


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
        prompt: make(chan string, 1),
        done: make(chan string, 1),
        doneReadCount: 0,
        promptStr: "",
        doneStr: "",
        readingPrompt: false,
        empty: false,

    }
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
}

func (o *OpenAIChat) chat(chat string, ctx context.Context) (string, error) {
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini20240718,
        Temperature: 0,
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
    row int
    col int
    ammo int
    level int
}

func NewTower(str string) Tower {
    assert.Assert(str[0] == '(' && str[len(str) - 1] == ')', "expected tower string to start and end with a paran", "str", str)

    parts := strings.Split(str[1:len(str) - 1], ",")
    assert.Assert(len(parts) == 4, "the tower definition did not contain 4 parts", "str", str)

    out := make([]int, 0, 4)
    for i, v := range parts {
        num, err := strconv.Atoi(v)
        assert.NoError(err, "unable to convert the number", "err", err)
        out[i] = num
    }

    return Tower{
        row: out[0],
        col: out[1],
        ammo: out[2],
        level: out[3],
    }
}

func NewTowers(str string) []Tower {
    out := []Tower{}
    start := 0
    for {
        start = strings.Index(str[start:], "(")
        if start == -1 {
            return out
        }

        end := strings.Index(str[start:], ")")
        assert.Assert(end != -1, "the tower string was incomplete", str)
        out = append(out, NewTower(str[start:start + end + 1]))

        start += end
    }
}

type GameState struct {
    rows string
    cols string
    allowedTowers int
    yourCreepDamage string
    enemyCreepDamage string
    yourTowers string
    enemyTowers string
    towerPlacement string
    creepSpawnRange string
}

func GameStateFromString(str string) GameState {
    lines := strings.Split(str, "\n")
    allowed, err := strconv.Atoi(strings.Split(lines[2], ": ")[1])
    assert.NoError(err, "unable to parse prompt allowed towers", "str", str)
    return GameState{
        rows: lines[0],
        cols: lines[1],
        allowedTowers: allowed,
        yourCreepDamage: lines[3],
        enemyCreepDamage: lines[4],
        yourTowers: lines[5],
        enemyTowers: lines[6],
        towerPlacement: lines[7],
        creepSpawnRange: lines[8],
    }
}

func (g *GameState) String() string {
    return fmt.Sprintf(`%s
%s
allowed towers: %d
%s
%s
%s
%s
%s
%s
`, g.rows, g.cols, g.allowedTowers, g.yourCreepDamage, g.enemyCreepDamage, g.yourTowers, g.enemyTowers, g.towerPlacement, g.creepSpawnRange)
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
            if viz {
                fmt.Printf("%s\n", string(b))
            }
            return len(b), nil
        })

    go cmdr.Run()
	go chtAgg.Pipe(twitchChat)

    ai := newStatefulOpenAIChat(os.Getenv("OPENAI_API_KEY"), string(systemPrompt), ctxWithCancel)
    done := make(chan struct{}, 1)
    aiRandomGuesses := 0
    aiBadParse := 0
    aiTotalTowers := 0
    roundCount := 0
    fmt.Printf("won,round,prompt file,seed,ai total towers,ai guesses,ai bad parses\n")
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
                fmt.Printf("1,%d,%s,%d,%d,%d,%d\n", roundCount, systemPromptFile, seed, aiTotalTowers, aiRandomGuesses, aiBadParse)
            } else {
                fmt.Printf("2,%d,%s,%d,%d,%d,%d\n", roundCount, systemPromptFile, seed, aiTotalTowers, aiRandomGuesses, aiBadParse)
            }

            done<-struct{}{}
        }
    }()

    defer func() {
        fmt.Println("\x1b[?25h")
    }()

    posGen := NewRandomPos(24)
    outer:
    for {
        debug.WriteStrLine(fmt.Sprintf("------------- waiting on game round: %d -----------------", roundCount))
        select {
        case <-ctxWithCancel.Done():
            break outer;
        case promptStr := <- errParser.prompt:
            prompt := GameStateFromString(promptStr)
            count := prompt.allowedTowers
            aiResponses := 0
            tries := 0
            roundCount++

            for aiResponses < count && tries < 3 {
                resp, err := ai.ReadWithTimeout(promptStr, time.Second * 5)
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

    <-done
}


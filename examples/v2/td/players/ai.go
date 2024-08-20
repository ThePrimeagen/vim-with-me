package players

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td/ai"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

func contextDone(ctx context.Context) bool {
    select {
    case <-ctx.Done():
        return true
    default:
        return false
    }
}

// if the ai provides some line that is unreadable i will just provide a value too large and thus random placement
func getPosFromAIResponse(line string) objects.Position {
    parts := strings.Split(line, ",")
    for i, p := range parts {
        parts[i] = strings.TrimSpace(p)
    }

    if len(parts) != 2 {
        return objects.OutOfBoundPosition()
    }

    row, err := strconv.Atoi(parts[0])
    if err != nil {
        return objects.OutOfBoundPosition()
    }

    col, err := strconv.Atoi(parts[1])
    if err != nil {
        return objects.OutOfBoundPosition()
    }
    return objects.Position{
        Row: uint(row),
        Col: uint(col),
    }
}

type AIResponder struct {
    ai AIFetcher
    maxTries uint
    debug *testies.DebugFile
    timeout time.Duration
    streamResults bool
    guesses int
    badParses int
    totalTowers int
}

type AIFetcher interface {
    ReadWithTimeout(prompt string, t time.Duration) (string, error)
}

func NewFetchPosition(ai AIFetcher, debug *testies.DebugFile) AIResponder {
    return AIResponder {
        ai: ai,
        debug: debug,
        maxTries: 3,
        timeout: time.Second * 5,
        streamResults: false,
        badParses: 0,
        guesses: 0,
        totalTowers: 0,
    };
}

func (f *AIResponder) Run(ctx context.Context) { }

func (f *AIResponder) Stats() objects.Stats {
    return objects.Stats{
        RandomGuesses: f.guesses,
        BadParses: f.badParses,
        TotalTowers: f.totalTowers,
    }
}

func (f *AIResponder) StartRound() {
    f.streamResults = true
    return
}

func (f *AIResponder) fetchResults(team uint8, gs *objects.GameState, ctx context.Context) []objects.Position {
    promptState := gs.PromptState(team)
    promptStr := promptState.String()
    f.debug.WriteStrLine(fmt.Sprintf("AIResponder#fetchResults Prompt(%d): \"%s\"", team, promptStr))
    resp, err := f.ai.ReadWithTimeout(promptState.String(), f.timeout)
    if contextDone(ctx) {
        return []objects.Position{}
    }

    if err != nil {
        f.debug.WriteStrLine(fmt.Sprintf("ai-placement error: \"%s\" err: \"%s\"", resp, err.Error()))
        parts := strings.Split(err.Error(), "try again in ")
        if len(parts) == 2 {
            secsStr := strings.Split(parts[1], " ")[0]
            secs, err := strconv.ParseFloat(secsStr[0:len(secsStr) - 2], 64)
            if err == nil {
                dur := time.Duration(float64(time.Second) * secs)
                f.debug.WriteStrLine(fmt.Sprintf("ai-placement wait time required: %d", dur))
                <-time.NewTimer(dur).C
            }
        }
        return []objects.Position{}
    }

    if resp == "" {
        <-time.NewTimer(time.Second).C
        return []objects.Position{}
    }

    responses := []objects.Position{}
    for _, line := range strings.Split(resp, "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            break;
        }

        parsedLine := getPosFromAIResponse(line)
        if parsedLine.OutOfBounds() {
            f.badParses++
            continue
        }

        responses = append(responses, parsedLine)
        if len(responses) >= gs.AllowedTowers {
            break
        }
    }

    f.debug.WriteStrLine(fmt.Sprintf("AIResponder#fetchResults Response(%d): %+v", team, responses))
    return responses
}

func (f *AIResponder) StreamResults(team uint8, gs *objects.GameState, out PositionChan, done Done, ctx context.Context) {
    if !f.streamResults {
        return
    }

    f.streamResults = false

    go func() {
        count := gs.AllowedTowers
        responses := []objects.Position{}
        var tries uint = 0

        for len(responses) < count && tries < f.maxTries {
            responses = append(responses, f.fetchResults(team, gs, ctx)...)
        }

        for range count - len(responses) {
            f.guesses++
            responses = append(responses, objects.OutOfBoundPosition())
        }

        if !contextDone(ctx) {
            out <- responses
            done <- struct{}{}
        }
    }()
}

func AIPlayerFromString(arg string, debug *testies.DebugFile, key string, ctx context.Context) AIResponder {
    assert.Assert(strings.HasPrefix(arg, "ai"), "invalid player string for ai client", "arg", arg)

    parts := strings.Split(arg, ":")
    assert.Assert(len(parts) == 3, "invalid ai player string colon count", "parts", parts)
    assert.Assert(parts[1] == "openai", "unsupported ai model for player", "parts", parts)

    systemPrompt, err := os.ReadFile(parts[2])
    assert.NoError(err, "could not open system prompt", "parts", parts)

    ai := ai.NewStatefulOpenAIChat(key, string(systemPrompt), ctx)
    return NewFetchPosition(ai, debug)
}

func GetPromptName(arg string) string {
    if !strings.HasPrefix(arg, "ai") {
        return ""
    }

    return strings.Split(arg, ":")[1]
}

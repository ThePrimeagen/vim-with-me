package players

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
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
    towersCreatedThisRound int
    team uint8
}

type AIFetcher interface {
    ReadWithTimeout(prompt string, t time.Duration) (string, error)
    Name() string
}

func NewFetchPosition(ai AIFetcher, team uint8, debug *testies.DebugFile) AIResponder {
    return AIResponder {
        ai: ai,
        debug: debug,
        maxTries: 1,
        timeout: time.Second * 5,
        streamResults: false,
        towersCreatedThisRound: 0,
        badParses: 0,
        guesses: 0,
        totalTowers: 0,
        team: team,
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
    f.towersCreatedThisRound = 0
    return
}

func isTimeoutError(e error, debug *testies.DebugFile) (bool, time.Duration) {
    // OpenAIs version
    err := e.Error()
    parts := strings.Split(err, "try again in ")
    if len(parts) == 2 {
        secsStr := strings.Split(parts[1], " ")[0]
        secs, err := strconv.ParseFloat(secsStr[0:len(secsStr) - 2], 64)
        if err == nil {
            dur := time.Duration(float64(time.Second) * secs)
            return true, dur
        }
    }


    // anthropics
    parts = strings.Split(err, "Number of request tokens has exceeded your per-minute rate limit")
    debug.WriteStrLine(fmt.Sprintf("ANTHROPIC TEST: %+v", parts))
    if len(parts) == 2 {
        return true, time.Second * 8
    }

    return false, time.Second
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
        if timeout, timeToWait := isTimeoutError(err, f.debug); timeout {
            f.debug.WriteStrLine(fmt.Sprintf("ai-placement timeout: \"%s\" for: \"%d\"", f.ai.Name(), timeToWait))
            timer := time.NewTimer(timeToWait)
            <-timer.C
            timer.Stop()
        }
        return []objects.Position{}
    }

    if resp == "" {
        return []objects.Position{}
    }

    responses := []objects.Position{}
    f.debug.WriteStrLine(fmt.Sprintf("AIResponder#fetchResults Raw Response(%d): \"%s\"", team, resp))
    for _, line := range strings.Split(resp, "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            continue;
        }

        parsedLine := getPosFromAIResponse(line)
        if parsedLine.OutOfBounds(gs) {
            f.badParses++
            continue
        }

        f.debug.WriteLine([]byte(fmt.Sprintf("parsedLine: %s", parsedLine.String())))
        responses = append(responses, parsedLine)
        if len(responses) >= gs.AllowedTowers {
            break
        }
    }

    if len(responses) < gs.AllowedTowers {
        for range gs.AllowedTowers - len(responses) {
            responses = append(responses, objects.OutOfBoundPosition())
        }
    }

    f.debug.WriteStrLine(fmt.Sprintf("AIResponder#fetchResults Response(%d): %+v", team, responses))
    return responses
}

func (f *AIResponder) Name() string {
    return f.ai.Name()
}

func (f *AIResponder) EndRound(gs *objects.GameState, cmdr td.TDCommander) {
    if f.towersCreatedThisRound >= gs.AllowedTowers {
        return
    }

    out := []objects.Position{}
    amount := gs.AllowedTowers - f.towersCreatedThisRound
    for range amount {
        out = append(out, objects.OutOfBoundPosition())
    }
    cmdr.WritePositions(out, f.team)
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

        for len(responses) < count && tries < f.maxTries && !contextDone(ctx) {
            responses = append(responses, f.fetchResults(team, gs, ctx)...)
        }

        for range count - len(responses) {
            f.guesses++
            responses = append(responses, objects.OutOfBoundPosition())
        }

        if !contextDone(ctx) {
            out <- responses
            f.towersCreatedThisRound = len(responses)
            done <- struct{}{}
        }
    }()
}

func AIPlayerFromString(arg string, team uint8, debug *testies.DebugFile, ctx context.Context) AIResponder {
    assert.Assert(strings.HasPrefix(arg, "ai"), "invalid player string for ai client", "arg", arg)

    parts := strings.Split(arg, ":")
    assert.Assert(len(parts) == 3, "invalid ai player string colon count", "parts", parts)

    systemPrompt, err := os.ReadFile(parts[2])
    assert.NoError(err, "could not open system prompt", "parts", parts)

    var fetcher AIFetcher = nil
    switch (parts[1]) {
    case "openai":
        fetcher = ai.NewStatefulOpenAIChat(string(systemPrompt), ctx)
    case "anthropic":
        fetcher = ai.NewClaudeSonnet(string(systemPrompt), ctx)
    }

    assert.Assert(fetcher != nil, "unsupported ai model for player", "parts", parts)
    return NewFetchPosition(fetcher, team, debug)
}

func GetPromptName(arg string) string {
    if !strings.HasPrefix(arg, "ai") {
        return ""
    }

    return strings.Split(arg, ":")[1]
}

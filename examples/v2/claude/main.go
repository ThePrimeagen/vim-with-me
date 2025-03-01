package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/ai"
)

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
}

func getResults(thinkingData AnthropicResponse) string {
	for _, content := range thinkingData.Content {
		if content.Type == "text" {
			return content.Text
		}
	}
	panic("WTF NO CONTENT???")
}

const SYSTEM = `
You are playing a tower defense game and will be communicated the state of the
game via json.

The json will contain the state of the game, the state of the enemy, and your
state.

Towers are described with 4 properties, row, col, level, ammo.

If your ammo reaches 0, the tower will be destroyed
If you have no ammo left, you will lose the game.

To place a tower specify in Row,Col coordinate system.  Each tower separated by
new lines.  You can place up to "allowedTowers" amount of towers (in json
prompt).  you must specify "allowedTowers" amount of R,C lines in your turn.

towers are 3x5 in size.  Do not build outside of the play area, specified in the game state's Rows and Cols

creep frequency will increase every round

Your goal is to win by surviving the longest.  You can use your towers to attack, to build new towers, or to upgrade existing towers

If you have a tower with low health, place a tower on the tower and it will upgrade it.
upgrades increase damage and range of towers.  Level 9 is OP and the final
level.  You cannot upgrade past Level 9.

once at level 9 you can still place a tower on it to replenish the ammo

You CANNOT place tower on COL 0.  Creeps spawn there and your tower will be randomly placed, which is bad

DO NOT EXPLAIN WHY YOU MAKE YOUR MOVES

Tower Upgrades:
Level 1: 1 damage, 1 range
Level 3: 1 damage, 2 range
Level 5: 2 damage, 3 range
Level 7: 2 damage, 4 range
Level 9: 3 damage, 6 range

Explaining Data Fields:
type Tower struct {
	row int
	col int
	ammo int
	level int
}

type Range struct {
    startRow uint
    endRow uint
}

type GameState struct {
    rows uint -- the number of rows in the game
    cols uint -- the number of columns in the game
    allowedTowers int -- the number of towers you can place this round
    yourCreepDamage uint -- the amount of damage you receive per tower when a creep gets to the end of your side
    enemyCreepDamage uint
    yourTowers []Tower -- your current towers row, col and ammo (its life) and level.  9 is max
    enemyTowers []Tower -- enemy towers
    yourTowerPlacementRange Range -- Where you can place your towers
    enemyTowerPlacementRange Range
    yourCreepSpawnRange Range -- the range of rows that the creeps you MUST kill will spawn
    enemyCreepSpawnRange Range
    round uint -- the current round
}

`

const PROMPT_STRING = "{\"rows\":24,\"cols\":80,\"allowedTowers\":1,\"yourCreepDamage\":1,\"enemyCreepDamage\":1,\"yourTowers\":[],\"TwoTowers\":[],\"yourTowerPlacementRange\":{\"startRow\":0,\"endRow\":4},\"enemyTowerPlacementRange\":{\"startRow\":20,\"endRow\":21},\"yourCreepSpawnRange\":{\"startRow\":0,\"endRow\":4},\"enemyCreepSpawnRange\":{\"startRow\":20,\"endRow\":24},\"round\":1}"

var msgs = []map[string]any{
	{"role": "user", "content": PROMPT_STRING},
}

func makeRequest() {
	payload, err := json.Marshal(map[string]interface{}{
		"model": "claude-3-7-sonnet-latest",
		"max_tokens": 8192,
		"system": SYSTEM,
		"messages": msgs,
	})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("x-api-key", os.Getenv("ANTHROPIC_API_KEY"))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("unable to read body data: %s", err))
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("non 200 status code: %d -- %s", resp.StatusCode, data))
	}

	defer resp.Body.Close()

	var result AnthropicResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	fmt.Printf("results: %s\n", getResults(result))
}

func makeOpenAIRequest() {
	ctx := context.Background()
	chat := ai.NewStatefulOpenAIChat(SYSTEM, ctx)

	out, err := chat.ReadWithTimeout(PROMPT_STRING, time.Second * 30)

	slog.Error("open ai response", "out", out, "err", err)
}

func main() {
	godotenv.Load()
	makeRequest()
	makeOpenAIRequest()
}


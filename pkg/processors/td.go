package processors

import (
	"strconv"
	"strings"
	"time"
)

type TDProcessor struct {
	out          chan string
	maxOccCount  int
	maxOcc       int
	pointChannel chan TDPoint
	points       map[int]int
}

type TDPoint struct {
	x int
	y int
}

func tdPointFromString(str string) *TDPoint {
	parts := strings.SplitN(str, ":", 3)
	if len(parts) != 3 {
		return nil
	}

	if parts[0] != "t" {
		return nil
	}

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	if x < 1 || y < 1 || x > 80 || y > 40 {
		return nil
	}

	return &TDPoint{x, y}
}

func NewTDProcessor(seconds int) *TDProcessor {

	processor := &TDProcessor{out: make(chan string, 10), pointChannel: make(chan TDPoint, 10), points: make(map[int]int)}

	go func() {
		ticker := time.NewTicker(time.Duration(seconds) * time.Second)
		for {
			select {

			case point := <-processor.pointChannel:
				p := point.y*1000 + point.x
				processor.points[p]++

				if processor.points[p] > processor.maxOccCount {
					processor.maxOcc = p
					processor.maxOccCount = processor.points[p]
				}

			case <-ticker.C:
				x := processor.maxOcc % 1000
				y := processor.maxOcc / 1000

				if x != 0 && y != 0 {
					processor.out <- "t:" + strconv.Itoa(x) + ":" + strconv.Itoa(y)
				}

				processor.points = make(map[int]int)
				processor.maxOcc = 0
				processor.maxOccCount = 0
			}
		}
	}()

	return processor
}

func (td *TDProcessor) Process(str string) {
	parts := strings.SplitN(str, ":", 3)
	if len(parts) != 3 {
		return
	}

	if parts[0] != "message" {
		return
	}

	contents := tdPointFromString(parts[2])
	if contents != nil {
		td.pointChannel <- *contents
	}
}

func (td *TDProcessor) Out() chan string {
	return td.out
}

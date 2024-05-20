package ascii_buffer

import (
	"fmt"
	"slices"

	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type Frequency struct {
	pointMap map[int]*FreqPoint
	count    int
	Points   []*FreqPoint
}

type FreqPoint struct {
	Val   int
	Count int
}

func NewFreqency() Frequency {
	return Frequency{
		pointMap: map[int]*FreqPoint{},
		Points:   make([]*FreqPoint, 0),
		count:    0,
	}
}

func (f *Frequency) Length() int {
	return f.count
}

func (f *Frequency) Reset() {
	f.count = 0
	f.pointMap = map[int]*FreqPoint{}
	f.Points = make([]*FreqPoint, 0)
}

func (f *Frequency) Freq(data byteutils.ByteIterator) {
	for {
		val := data.Next()
		point, ok := f.pointMap[val.Value]

		if !ok {
			f.count++

			point = &FreqPoint{
				Count: 0,
				Val:   val.Value,
			}
			f.Points = append(f.Points, point)
			f.pointMap[val.Value] = point
		}

		point.Count++

		// i wish i had a do while... the lords loop
		if val.Done {
			break
		}
	}
}

func (f *Frequency) DebugFunc(toString func(int) string) string {
	points := make([]*FreqPoint, len(f.Points), len(f.Points))
	copy(points, f.Points)
	slices.SortFunc(points, func(a, b *FreqPoint) int {
		return a.Count - b.Count
	})

	out := fmt.Sprintf("Frequency(%d): ", len(points))
	top := 0
	bottom := 0
	for i, p := range points {
		if i+4 >= len(points) {
			top += p.Count
		} else {
			bottom += p.Count
		}
		out += fmt.Sprintf("%s(%d) ", toString(p.Val), p.Count)
	}

	out += fmt.Sprintf("-- top %d bottom %d diff %d", top, bottom, top-bottom)

	return out
}

func (f *Frequency) Debug() string {
	return f.DebugFunc(func(b int) string {
		return fmt.Sprintf("%d", b)
	})
}

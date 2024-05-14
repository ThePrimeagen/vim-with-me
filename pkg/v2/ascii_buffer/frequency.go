package ascii_buffer

import (
	"fmt"
	"slices"

	"github.com/theprimeagen/vim-with-me/pkg/v2/iterator"
)

type Frequency struct {
	points map[int]*FreqPoint
	count  int
	Points []*FreqPoint
}

type FreqPoint struct {
	val   int
	count int
}

func NewFreqency() Frequency {
	return Frequency{
		points: map[int]*FreqPoint{},
		Points: make([]*FreqPoint, 0),
		count:  0,
	}
}

func (f *Frequency) Length() int {
	return f.count
}

func (f *Frequency) Reset() {
	f.count = 0
    f.points = map[int]*FreqPoint{}
    f.Points = make([]*FreqPoint, 0)
}

func (f *Frequency) Freq(data iterator.ByteIterator) {
	for {
        val := data.Next()
        point, ok := f.points[val.Value]

		if !ok {
			f.count++

			point = &FreqPoint{count: 0}
			f.Points = append(f.Points, point)
		}

		point.count++

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
		return a.count - b.count
	})

	out := fmt.Sprintf("Frequency(%d): ", len(points))
	top := 0
	bottom := 0
	for i, p := range points {
		if i+4 >= len(points) {
			top += p.count
		} else {
			bottom += p.count
		}
		out += fmt.Sprintf("%s(%d) ", toString(p.val), p.count)
	}

	out += fmt.Sprintf("-- top %d bottom %d diff %d", top, bottom, top-bottom)

	return out
}

func (f *Frequency) Debug() string {
	return f.DebugFunc(func(b int) string { return string(b) })
}

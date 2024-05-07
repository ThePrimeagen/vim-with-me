package ascii_buffer

import (
	"fmt"
	"slices"
)

type Frequency struct {
	Buffer [256]int
    count int
    Points []*FreqPoint
    possiblePoints [256]*FreqPoint
}

type FreqPoint struct {
	idx   byte
	count int
}

func NewFreqency() Frequency {
	return Frequency{
		Buffer: [256]int{},
		possiblePoints: [256]*FreqPoint{},
        Points: make([]*FreqPoint, 0),
        count: 0,
	}
}

func (f *Frequency) Length() int {
    return f.count
}

func (f *Frequency) Reset() {
    f.count = 0
	for i := 0; i < 256; i++ {
		f.Buffer[i] = 0
	}
}

func (f *Frequency) Freq(data []byte) {
	for _, b := range data {
        if f.Buffer[b] == 0 {
            f.count++

            point := &FreqPoint{count: 0, idx: b}
            f.possiblePoints[b] = point
            f.Points = append(f.Points, point)
        }

		f.Buffer[b]++
        f.possiblePoints[b].count++
	}
}

func (f *Frequency) Debug() string {
    points := make([]*FreqPoint, len(f.Points), len(f.Points))
    copy(points, f.Points)
    slices.SortFunc(points, func(a, b *FreqPoint) int {
        return a.count - b.count
    })

    out := fmt.Sprintf("Frequency(%d): ", len(points))
    top := 0
    bottom := 0
    for i, p := range points {
        if i + 4 >= len(points) {
            top += p.count
        } else {
            bottom += p.count
        }
        out += fmt.Sprintf("%s(%d) ", string(p.idx), p.count)
    }

    out += fmt.Sprintf("-- top %d bottom %d diff %d", top, bottom, top - bottom)

    return out
}

package ascii_buffer

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

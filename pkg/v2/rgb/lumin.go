package rgb

type BrightnessRange struct {
	Max  float64
	Min  float64
	Avg  float64
	Char string
}

type Lumin struct {
	ranges []BrightnessRange
	buffer []byte
	idx    int
}

var defaultBrightnessRange = []BrightnessRange{}

func NewLumin(count int) *Lumin {
	return &Lumin{
		buffer: make([]byte, count, count),
		ranges: defaultBrightnessRange,
		idx:    0,
	}
}

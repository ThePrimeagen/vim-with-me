package ascii_buffer

type Frequency struct {
    Buffer [256]int
}

func NewFreqency() Frequency {
    return Frequency{
        Buffer: [256]int{},
    }
}

func (f *Frequency) Reset() {
    for i := 0; i < 256; i++ {
        f.Buffer[i] = 0
    }
}

func (f *Frequency) Freq(data []byte) {
    for _, b := range data {
        f.Buffer[b]++
    }
}

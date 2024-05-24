package chat

import "fmt"

type ChatMsg struct {
	Name string
	Msg  string
	Bits int
}

func (c *ChatMsg) String() string {
    if c.Bits > 0 {
        return fmt.Sprintf("ChatMsg(bits: %d from: %s): %s", c.Bits, c.Name, c.Msg)
    }
    return fmt.Sprintf("ChatMsg(%s): %s", c.Name, c.Msg)
}

type FilterCB func(msg string) bool
type MapCB func(msg string) string
type Occurrence struct {
	Count int
	Msg   string
}

type ChatAggregator struct {
	mapFn       MapCB
	filter      FilterCB
	occurrences []*Occurrence
	max         *Occurrence
}

func (c *Occurrence) String() string {
    return fmt.Sprintf("Occurrence(%d): %s", c.Count, c.Msg)
}

var defaultMax = Occurrence{Count: 0, Msg: ""}
var identity = func(x string) string { return x }

func NewChatAggregator() ChatAggregator {
	return ChatAggregator{
		filter:      nil,
		occurrences: make([]*Occurrence, 0),
		max:         &defaultMax,
		mapFn:       identity,
	}
}

func (c ChatAggregator) WithFilter(filterFn FilterCB) ChatAggregator {
	c.filter = filterFn
	return c
}

func (c ChatAggregator) WithMap(mapFn MapCB) ChatAggregator {
	c.mapFn = mapFn
	return c
}

func (c *ChatAggregator) incAndMax(occurrence *Occurrence) {
    occurrence.Count++
    if occurrence.Count > c.max.Count {
        c.max = occurrence
    }
}

func (c *ChatAggregator) Add(msg string) {
	msg = c.mapFn(msg)
	if !c.filter(msg) {
        fmt.Printf("filtering msg: %s\n", msg)
		return
	}

	for _, occurrence := range c.occurrences {
		if occurrence.Msg == msg {
            c.incAndMax(occurrence)
			return
		}
	}

    occurrence := &Occurrence{Count: 0, Msg: msg}
	c.occurrences = append(c.occurrences, occurrence)
    c.incAndMax(occurrence)
}

func (c *ChatAggregator) String() string {
    out := fmt.Sprintf("ChatAggregator:(%v)\n", c.max)

    for _, o := range c.occurrences {
        out += fmt.Sprintf("%v\n", o)
    }

    return out
}

func (c *ChatAggregator) Reset() Occurrence {
    fmt.Printf("%v\n", c)
	max := c.max
	c.max = &defaultMax
	c.occurrences = make([]*Occurrence, 0)

	return *max
}

func (c *ChatAggregator) Next() string {
    return c.Reset().Msg
}

func (c *ChatAggregator) Pipe(ch chan ChatMsg) {
	for msg := range ch {
        fmt.Printf("chat: %v\n", msg)
		c.Add(msg.Msg)
	}
}

package main

import (
	"fmt"
)

type Color struct {
    Foo int
}

type Cell struct {
	Foreground Color
	Background Color
	Value      byte
}

func main() {

    c := []Cell{Cell{
        Foreground: Color{Foo: 1},
        Background: Color{Foo: 10},
        Value: 100,
    }}

    d := c[0]
    d.Value += 1

    fmt.Printf("values: %+v %+v\n", c, d)
}

package main

import (
	"context"
	"fmt"
	"time"
)
func one(ctx context.Context, ch chan string) {
    go func() {
        Outer:
        for {
            select {
            case <-ctx.Done():
                foo := ctx.Value("foo")
                fmt.Printf("done with context, one: %+v\n", foo)
                break Outer;
            case ch <- "test":
                fmt.Printf("got a value from channel")
            }
        }
    }()
}



func two(ctx context.Context, ch chan string) {
    ctx = context.WithValue(ctx, "foo", "bar")
    one(ctx, ch)
    go func() {
        Outer:
        for {
            select {
            case <-ctx.Done():
                fmt.Printf("done with context, two\n")
                break Outer;
            case ch <- "test":
                fmt.Printf("got a value from channel")
            }
        }
    }()
}


func three(ctx context.Context, ch chan string) {
    go func() {
        Outer:
        for {
            select {
            case <-ctx.Done():
                fmt.Printf("done with context, three\n")
                break Outer;
            case ch <- "test":
                fmt.Printf("got a value from channel")
            }
        }
    }()
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    ch := make(chan string)
    one(ctx, ch)
    two(ctx, ch)
    three(ctx, ch)

    <-time.NewTimer(3 * time.Second).C
    cancel()
    <-time.NewTimer(1 * time.Second).C
}

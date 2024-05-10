package assert

import (
	"fmt"
	"log"
	"log/slog"
)

var assertData map[string]any = map[string]any{}
func AddAssertData(key string, value any) {
    assertData[key] = value
}

func RemoveAssertData(key string) {
    delete(assertData, key)
}

func runAssert(msg string, args ...any) {
    for k, v := range assertData {
        slog.Error("context", "key", k, "value", v)
    }
    fmt.Printf("%s: ", msg)
    for _, item := range args {
        fmt.Printf("%+v ", item)
    }
    log.Fatal("runtime assert failure")
}

// TODO: Think about passing around a context for debugging purposes
func Assert(truth bool, msg string, data ...any) {
    if !truth {
        runAssert(msg, data...)
    }
}

func NoError(err error, msg string) {
    if err != nil {
        slog.Error("NoError#error encountered", "error", err)
        runAssert(msg)
    }
}


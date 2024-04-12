package assert

import (
	"log"
)

// TODO: Think about passing around a context for debugging purposes
func Assert(truth bool, msg string) {
    if !truth {
        log.Fatal(msg)
    }
}


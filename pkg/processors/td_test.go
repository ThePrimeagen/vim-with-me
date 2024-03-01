package processors_test

import (
	"testing"

	"chat.theprimeagen.com/pkg/processors"
)

func TestTDProcessor(t *testing.T) {
    td := processors.NewTDProcessor(1)

    td.Process("message:shutupchat:t:1:1")
    td.Process("message:shutupchat:t:1:1")
    td.Process("message:shutupchat:t:2:2")
    td.Process("message:shutupchat:t:69:20")
    td.Process("message:shutupchat:t:69:20")
    td.Process("message:shutupchat:t:69:20")

    res := <-td.Out()

    if res != "t:69:20" {
        t.Errorf("Expected t:69:20, got %s", res)
    }
}


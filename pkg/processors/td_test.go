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

    td.Process("message:shutupchat:t:2:3")
    td.Process("message:shutupchat:t:6:9")
    td.Process("message:shutupchat:t:4:20")
    td.Process("message:shutupchat:t:20:20")
    td.Process("message:shutupchat:t:6:9")
    td.Process("message:shutupchat:t:6:9")
    td.Process("message:shutupchat:t:6:9")
    td.Process("message:shutupchat:t:4:20")

    res = <-td.Out()

    if res != "t:6:9" {
        t.Errorf("Expected t:6:9, got %s", res)
    }


    td.Process("message:bearow:t:21:37")
    td.Process("message:mSyke:t:10:10")
    td.Process("message:Rub1_NN:t:4:20")
    td.Process("message:sedgrepsupercombo:t:40:30")
    td.Process("message:FV7VR3:t:20:18")
    td.Process("message:RegardedGlizzy:t:6:9")
    td.Process("message:imastonedhippy:t:１0:１3")
    td.Process("message:FV7VR3:t:20:1")
    td.Process("message:sonicfind:t:4:20")
    td.Process("message:FV7VR3:t:20:18")

    res = <-td.Out()

    if res != "t:4:20" {
        t.Errorf("Expected t:4:20, got %s", res)
    }


}

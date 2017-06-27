package test

import (
	"testing"

	"github.com/ariefdarmawan/keepalive"
	"github.com/eaciit/toolkit"
)

func TestParse(t *testing.T) {
	cmd := "keepalive -config=\"d:\\users\\ariefdarmawan\\some thing\""
	cmds := keepalive.CmdToStrings(cmd)
	if len(cmds) != 2 {
		t.Errorf("Got: %v len is: %d", toolkit.JsonString(cmds), len(cmds))
	} else {
		toolkit.Printf("Got: %v \n", toolkit.JsonString(cmds))
	}
}

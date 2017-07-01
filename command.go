package keepalive

import (
	"bytes"
	"os/exec"

	"io/ioutil"

	"github.com/eaciit/toolkit"
)

type ServiceCommandTypeEnum string

const (
	ServiceCheck   ServiceCommandTypeEnum = "Check"
	ServiceStart                          = "Start"
	ServiceStop                           = "Stop"
	ServiceRestart                        = "Restart"
)

type OpEnum string

const (
	OpEq          OpEnum = "eq"
	OpNeq                = "ne"
	OpContains           = "contains"
	OpNotContains        = "notcontains"
)

type Command struct {
	CommandType CommandTypeEnum
	//Function    ServiceCommandTypeEnum
	Txt string

	Op            OpEnum
	Expected      string
	CaseSensitive bool
}

func (c *Command) Exec() *toolkit.Result {
	if c.CommandType == CommandUrl {
		return runRest(c.Txt, "GET", nil)
	} else if c.CommandType == CommandLine {
		return runCmd(CmdToStrings(c.Txt)...)
	}

	return toolkit.NewResult().SetErrorTxt("Not implemented")
}

func CmdToStrings(txt string) []string {
	var ret []string
	splitchar := "\""
	delim := " "
	hasprefix := false
	tmp := ""
	for _, c := range txt {
		cs := string(c)
		if cs != delim && cs != splitchar {
			tmp += cs
		} else if cs == splitchar {
			tmp += cs
			hasprefix = !hasprefix
		} else if cs == delim && hasprefix {
			tmp += cs
		} else if cs == delim && !hasprefix {
			ret = append(ret, tmp)
			tmp = ""
		}
	}
	if tmp != "" {
		ret = append(ret, tmp)
	}
	return ret
}

func runCmd(cmds ...string) *toolkit.Result {
	ret := toolkit.NewResult()

	var (
		cmd            *exec.Cmd
		outStr, errStr string
	)
	if len(cmds) >= 2 {
		cmd = exec.Command(cmds[0], cmds[1:]...)
	} else {
		cmd = exec.Command(cmds[0])
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		//return ret.SetErrorTxt(fmt.Sprintf("Unable to exec %s: %s", cmds[0], err.Error()))
		return ret.SetError(err)
	}
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())

	if errStr != "" {
		return ret.SetErrorTxt(errStr)
		//return ret.SetErrorTxt(fmt.Sprintf("Unable to exec %s: %s", cmds[0], errStr))
	}

	ret.SetData(outStr)
	return ret
}

func runRest(url, method string, data []byte) *toolkit.Result {
	ret := toolkit.NewResult()
	r, err := toolkit.HttpCall(url, method, data, nil)
	if err != nil {
		return ret.SetErrorTxt(toolkit.Sprintf("Unable to call %s: %s", url, err.Error()))
	}

	defer r.Body.Close()

	output, eread := ioutil.ReadAll(r.Body)
	if eread != nil {
		return ret.SetErrorTxt(toolkit.Sprintf("Unable to read %s: %s", url, eread.Error()))
	}
	return ret.SetData(string(output))
}

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ariefdarmawan/keepalive"
	"github.com/eaciit/toolkit"
)

var (
	fconfig = flag.String("config", "", "-config=\"path to config file\" define config file, if blank it will be assume on same location of the executable with name keepalive.json")
	log, _  = toolkit.NewLog(true, false, "", "", "")

	ctx keepalive.Context
)

func main() {
	flag.Parse()
	loadconfig()

	if e := ctx.Run(); e != nil {
		return
	}
	ctx.Wait()
}

func loadconfig() {
	configpath := *fconfig
	if configpath == "" {
		configpath, _ = os.Executable()
		configpath = filepath.Join(filepath.Dir(configpath), "keepalive.conf")
	}

	bytepaths, err := ioutil.ReadFile(configpath)
	if err != nil {
		log.Error(toolkit.Sprintf("Reading %s fail: %s", configpath, err.Error()))
		os.Exit(10)
	}

	var configValue map[string]interface{}
	err = json.Unmarshal(bytepaths, &configValue)
	if err != nil {
		log.Error(toolkit.Sprintf("Unmarshalling %s fail: %s", configpath, err.Error()))
		os.Exit(10)
	}

	err = toolkit.Serde(configValue, &ctx, "")
	if err != nil {
		log.Error(toolkit.Sprintf("Serde %s fail: %s", configpath, err.Error()))
		os.Exit(10)
	}
}

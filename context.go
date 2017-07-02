package keepalive

import (
	"os"
	"path/filepath"

	"github.com/eaciit/toolkit"
)
import "sync"

type CommandTypeEnum string

const (
	CommandLine CommandTypeEnum = "Cmd"
	CommandUrl                  = "Url"
)

type Context struct {
	Port       int
	Verbose    bool
	Services   map[string]*Service
	SmtpClient *SmtpClient
	LogToFile  bool
	LogPath    string

	ch  chan bool
	log *toolkit.LogEngine
}

func (c *Context) Run() error {
	if c.log == nil {
		c.log, _ = toolkit.NewLog(true, true, c.LogPath, "ctx-%s.log", "yyyy-MM-dd")
	}
	c.ch = make(chan bool)

	c.log.Info("Initiating context")
	c.StartWebConsole()

	activeServiceCount := 0
	createFolderIfNotExist(c.LogPath)
	serviceLogPath := filepath.Join(c.LogPath, "services")
	createFolderIfNotExist(serviceLogPath)
	for n, s := range c.Services {
		s.name = n
		s.verbose = c.Verbose
		s.smtpClient = c.SmtpClient

		if s.Active {
			s.log, _ = toolkit.NewLog(true, true, serviceLogPath, s.name+"-%s.log", "yyyy-MM-dd")
			s.StartMonitor()
			activeServiceCount++
		}
	}

	if activeServiceCount == 0 {
		c.log.Infof("No active service is defined. System should be stopped")
		return toolkit.Error("No active service is defined. System should be stopped")
	}

	return nil
}

func (c *Context) Wait() {
	for {
		select {
		case stop := <-c.ch:
			if stop {
				closewg := new(sync.WaitGroup)

				c.log.Info("KeepAlive will be stopped")
				for _, s := range c.Services {
					if s.Active {
						closewg.Add(1)
						s.closewg = closewg
						s.chStop <- true
					}
				}

				closewg.Wait()
				c.log.Info("KeepAlive is stopped")
				return
			}
		}
	}
}

func (c *Context) Stop() {
	c.ch <- true
}

func createFolderIfNotExist(path string) {
	if s, e := os.Stat(path); e == nil && s.IsDir() {
		return
	}

	os.Mkdir(path, os.FileMode(0644))
}

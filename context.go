package keepalive

import "github.com/eaciit/toolkit"
import "sync"

type CommandTypeEnum int

const (
	CommandLine CommandTypeEnum = 0
	CommandUrl                  = 1
)

type Context struct {
	Port     int
	Verbose  bool
	Services map[string]*Service

	ch  chan bool
	log *toolkit.LogEngine
}

func (c *Context) Run() error {
	if c.log == nil {
		c.log, _ = toolkit.NewLog(true, false, "", "", "")
	}
	c.ch = make(chan bool)

	c.log.Info("Initiating context")
	c.StartWebConsole()

	activeServiceCount := 0
	for n, s := range c.Services {
		s.log = c.log

		s.name = n
		s.verbose = c.Verbose

		if s.Active {
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

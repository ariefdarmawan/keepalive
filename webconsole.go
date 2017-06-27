package keepalive

import (
	"net/http"

	"time"

	"github.com/eaciit/toolkit"
)

func (c *Context) StartWebConsole() {
	c.log.Infof("Initiating web console on port %d", c.Port)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("KeepAlive (c) Arief Darmawan"))
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Daemon is stopped\n"))
		time.Sleep(1 * time.Second)
		c.ch <- true
	})

	go func() {
		serverurl := toolkit.Sprintf(":%d", c.Port)
		http.ListenAndServe(serverurl, nil)
	}()
}

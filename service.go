package keepalive

import (
	"time"

	"sync"

	"strings"

	"github.com/eaciit/toolkit"
)

type ServiceNotifEnum string

const (
	ServiceNotifStop    ServiceNotifEnum = "Service Stop"
	ServiceNotifRestart                  = "Service Restart"
	ServiceNotifFail                     = "Service Fail"
)

type Service struct {
	Commands        map[string]Command
	Interval        time.Duration
	NotifyInterval  time.Duration
	CheckAgainAfter time.Duration
	StopNotifyAfter int
	StopStartAfter  int
	Active          bool

	name      string
	verbose   bool
	chStop    chan bool
	lastCheck time.Time
	log       *toolkit.LogEngine
	closewg   *sync.WaitGroup

	startAttempt int
}

func (s *Service) run() {
	for {
		select {
		case stop := <-s.chStop:
			if stop {
				s.log.Infof("Service %s has been dismissed from to be monitored", s.name)
				if s.closewg != nil {
					s.closewg.Done()
				}
				return
			}

		case <-time.After(s.Interval * time.Millisecond):
			s.check()
		}
	}
}

func (s *Service) StartMonitor() {
	s.chStop = make(chan bool)
	if s.log == nil {
		s.log, _ = toolkit.NewLog(true, false, "", "", "")
	}
	go s.run()
	s.log.Infof("Service %s has been started to be monitored", s.name)
}

func (s *Service) StopMonitor() {
	s.chStop <- true
}

func (s *Service) check() {
	if s.StopStartAfter > 0 && s.startAttempt == s.StopStartAfter {
		time.Sleep(s.CheckAgainAfter * time.Millisecond)
		s.startAttempt = 0
	}

	checkErrMsg := ""
	r := s.Exec(ServiceCheck)
	if r.Status != toolkit.Status_OK {
		checkErrMsg = r.Message
	} else {
		cmd, _ := s.Commands[string(ServiceCheck)]
		if evalid := validate(r.Data.(string), cmd.Op, cmd.Expected, cmd.CaseSensitive); evalid != nil {
			checkErrMsg = evalid.Error()
		}
	}

	if checkErrMsg != "" {
		s.sendNotification(ServiceNotifStop, toolkit.M{}.Set("Message", checkErrMsg))
		if estart := s.restartService(); estart != nil {
			s.sendNotification(ServiceNotifFail, toolkit.M{}.Set("Message",
				toolkit.Sprintf("Restart %d service %s fail: %s", s.startAttempt, s.name, estart.Error())))
		}
	}

	if s.verbose {
		s.log.Infof("Service %s is OK", s.name)
	}
}

func (s *Service) restartService() error {
	if s.StopStartAfter > 0 && s.startAttempt == s.StopStartAfter {
		return nil
	}

	s.startAttempt++
	r := s.Exec(ServiceStart)
	if r.Status != toolkit.Status_OK {
		return toolkit.Error(r.Message)
	}

	cmd, _ := s.Commands[string(ServiceStart)]
	if evalid := validate(r.Data.(string), cmd.Op, cmd.Expected, cmd.CaseSensitive); evalid != nil {
		return evalid
	}
	s.sendNotification(ServiceNotifRestart, toolkit.M{}.Set("Message", toolkit.Sprintf("%s OK", s.name)))
	return nil
}

func (s *Service) sendNotification(notittype ServiceNotifEnum, in toolkit.M) {
	if in == nil {
		in = toolkit.M{}
	}

	msg := ""

	msg = in.GetString("Message")
	if notittype == ServiceNotifStop {
		s.log.Warningf("%s %s: %s", s.name, string(notittype), msg)
	} else if notittype == ServiceNotifFail {
		s.log.Errorf("%s: %s", string(notittype), msg)
	} else {
		s.log.Infof("%s: %s", string(notittype), msg)
	}
}

func (s *Service) Exec(ct ServiceCommandTypeEnum) *toolkit.Result {
	r := toolkit.NewResult()
	cmd, cmdfound := s.Commands[string(ct)]
	if !cmdfound {
		return r.SetErrorTxt(toolkit.Sprintf("Command %s could not be found", cmd))
	}

	r = cmd.Exec()
	return r
}

func validate(data string, op OpEnum, expected string, casesensitive bool) error {
	var datalow, expectedlow string
	if casesensitive {
		datalow = data
		expectedlow = expected
	} else {
		datalow = strings.ToLower(data)
		expectedlow = strings.ToLower(expected)
	}

	if op == OpEq && datalow == expectedlow {
		return nil
	} else if op == OpNeq && datalow != expectedlow {
		return nil
	} else if op == OpContains && strings.Index(datalow, expectedlow) >= 0 {
		return nil
	} else if op == OpNotContains && strings.Index(datalow, expectedlow) < 0 {
		return nil
	}

	return toolkit.Errorf("Invalid output, expecting output %s %s", string(op), expected)
}

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

func TestSendEmail(t *testing.T) {
	s := keepalive.SmtpClient{
		Host:     "smtp.office365.com",
		Port:     587,
		UserId:   "mailer@eaciit.com",
		Password: "Ruka6309",
		SSL:      false,
		TLS:      false,
	}

	msg := &keepalive.EmailMsg{
		From:    "mailer@eaciit.com",
		To:      []string{"arief@eaciit.com"},
		Cc:      []string{"ariefda@hotmail.com", "adarmawan.2006@gmail.com"},
		Subject: "Send email dari KeepAlive",
		Body: "Hi\n" +
			"Kindly please find status of KeepAlive: OK\n\nKeepAlive's team",
	}

	if err := s.Send(msg); err != nil {
		t.Errorf("Fail send email: %s", err.Error())
	}
}

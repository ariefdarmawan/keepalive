package keepalive

import (
	"crypto/tls"
	//"eaciit/gomail"

	"gopkg.in/gomail.v2"
)

type SmtpClient struct {
	Host, UserId, Password string
	Port                   int
	SSL                    bool
	TLS                    bool
}

type EmailMsg struct {
	From          string
	To, Cc, Bcc   []string
	Subject, Body string
	Attachments   []string
	HtmlBody      bool
}

func (s *SmtpClient) Send(msgs ...*EmailMsg) error {
	// Dialer
	conf := gomail.NewDialer(s.Host, s.Port, s.UserId, s.Password)
	if s.SSL {
		conf.SSL = s.SSL
	}
	if s.TLS {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Host,
		}
		conf.TLSConfig = tlsconfig
	}
	//conf.LocalName = s.Host

	// Dialer
	d, err := conf.Dial()
	if err != nil {
		return err
	}
	defer d.Close()

	// Message
	var ms []*gomail.Message
	for _, msg := range msgs {
		m := gomail.NewMessage()
		// Setup headers
		m.SetHeader("From", msg.From)
		m.SetHeader("To", msg.To...)
		if len(msg.Bcc) > 0 {
			m.SetHeader("Bcc", msg.Bcc...)
		}
		if len(msg.Cc) > 0 {
			m.SetHeader("Cc", msg.Cc...)
		}
		m.SetHeader("Subject", msg.Subject)

		body := msg.Body + "\n"
		if !msg.HtmlBody {
			//body += toolkit.Sprintf("Send timestamp: %s", toolkit.Date2String(time.Now(), "dd-MMM-yyyy hh:mm:ss"))
			m.SetBody("text/plain", body)
		} else {
			//body := msg.Body + "\n<br/>"
			//body += toolkit.Sprintf("<p>Send timestamp: %s</p>", toolkit.Date2String(time.Now(), "dd-MMM-yyyy hh:mm:ss"))
			m.SetBody("text/html", body)
		}

		if len(msg.Attachments) > 0 {
			for _, v := range msg.Attachments {
				m.Attach(v)
			}
		}
		ms = append(ms, m)
	}

	return gomail.Send(d, ms...)
}

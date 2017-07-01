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
			//InsecureSkipVerify: true,
			ServerName: s.Host,
		}
		conf.TLSConfig = tlsconfig
	}
	conf.LocalName = s.Host

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

		if !msg.HtmlBody {
			m.SetBody("text/plain", msg.Body)
		} else {
			m.SetBody("text/html", msg.Body)
		}

		if len(msg.Attachments) > 0 {
			for _, v := range msg.Attachments {
				m.Attach(v)
			}
		}
		ms = append(ms, m)
	}

	err = gomail.Send(d, ms[0])
	ms[0].Reset()
	return err
}

func (s *SmtpClient) send_(msgs ...*EmailMsg) error {
	// Dialer Config
	conf := gomail.NewDialer(s.Host, s.Port, s.UserId, s.Password)
	if s.SSL {
		conf.SSL = s.SSL
	}
	if s.TLS {
		// TLS config
		tlsconfig := &tls.Config{
			//InsecureSkipVerify: true,
			ServerName: s.Host,
		}
		conf.TLSConfig = tlsconfig
	}
	conf.LocalName = s.Host

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

		if !msg.HtmlBody {
			m.SetBody("text/plain", msg.Body)
		} else {
			m.SetBody("text/html", msg.Body)
		}

		if len(msg.Attachments) > 0 {
			for _, v := range msg.Attachments {
				m.Attach(v)
			}
		}
		ms = append(ms, m)
	}

	err = gomail.Send(d, ms[0])
	return err
}

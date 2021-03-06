package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/voioc/coco/logzap"
)

type PEmailConfig struct {
	Host string
	Port string
	User string
	Pwd  string
	From string
}

type PEmail struct {
	Config  *PEmailConfig
	Subject string
	Body    string
	To      string
	Format  string
}

var (
	mailChan chan *PEmail
	config   *PEmailConfig
)

func init() {
	poolSize := 10 // viper.GetInt("email.pool")
	host := viper.GetString("email.host")
	port := viper.GetString("email.port")
	user := viper.GetString("email.user")
	pwd := viper.GetString("email.password")
	from := viper.GetString("email.from")

	config = &PEmailConfig{
		Host: host,
		From: from,
		Port: port,
		User: user,
		Pwd:  pwd,
	}

	// 创建通道
	mailChan = make(chan *PEmail, poolSize)

	go func() {
		for {
			select {
			case m, ok := <-mailChan:
				if !ok {
					return
				}
				if err := m.SendToEmail(); err != nil {
					logzap.Ex(context.Background(), "SendMail:", err.Error())
				}
			}
		}
	}()
}

func SendToChan(to, subject, body, mailtype string) bool {
	email := &PEmail{
		Config:  config,
		Body:    body,
		Subject: subject,
		Format:  mailtype,
		To:      to,
	}
	select {
	case mailChan <- email:
		return true
	case <-time.After(time.Second * 3):
		return false
	}
}

func (pe *PEmail) SendToEmail() error {
	auth := smtp.PlainAuth("", pe.Config.User, pe.Config.Pwd, pe.Config.Host)
	contentType := GetContentTypeString(pe.Format)

	msg := []byte("To: " + pe.To + "\r\nFrom: " + pe.Config.User +
		"\r\nSubject: " + pe.Subject + "\r\n" + contentType + "\r\n\r\n" + pe.Body)

	sendTo := strings.Split(pe.To, ";")

	var err error

	if pe.Config.Port == "25" {
		err = SendMailUsing25(pe.Config.Host, pe.Config.Port, auth, pe.Config.User, sendTo, msg)
	} else if pe.Config.Port == "465" {
		err = SendMailUsing465(pe.Config.Host, pe.Config.Port, auth, pe.Config.User, sendTo, msg)
	} else {
		err = fmt.Errorf("%s", "other ports are not supported,please check the app.conf configuration file")
	}

	return err
}

func GetContentTypeString(format string) string {
	var contentType string
	if format == "" {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/" + format + "; charset=UTF-8"
	}
	return contentType
}

func SendMailUsing25(addr string, port string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	err = smtp.SendMail(addr+":"+port, auth, from, to, msg)
	return err
}

func SendMailUsing465(addr string, port string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	c, err := Dial(addr + ":" + port)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			fmt.Print(err)
			return err

		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

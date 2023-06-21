package email

import (
	"blog/pkg/l"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"net/smtp"
	"strings"
	"time"
)

type Email struct {
	Auth     smtp.Auth
	Host     string
	Port     string
	From     string
	Password string
	To       string
	Title    string
	Data     []byte
}

func SendMail(mail *Email) error {
	mail.Auth = smtp.PlainAuth("", mail.From, mail.Password, mail.Host)

	conn, err := tls.Dial("tcp", mail.Host, nil)
	if err != nil {
		l.Logger().Error("sendMail Dial failed", zap.Error(err))
		return err
	}
	cli, err := smtp.NewClient(conn, mail.Host)
	if err != nil {
		l.Logger().Error("sendMail get client failed", zap.Error(err))
		return err
	}
	defer cli.Close()

	if ok, _ := cli.Extension("AUTH"); ok && mail.Auth != nil {
		// 服务器如果支持AUTH扩展，则进行校验
		if err = cli.Auth(mail.Auth); err != nil {
			l.Logger().Error("sendMail Auth failed", zap.Error(err))
			return err
		}
	}

	if err = cli.Mail(mail.From); err != nil {
		l.Logger().Error("sendMail send MAIL command failed", zap.Error(err))
		return err
	}

	sendList := strings.Split(mail.To, ",")
	for _, addr := range sendList {
		if err = cli.Rcpt(addr); err != nil {
			l.Logger().Error("sendMail send RCPT command failed", zap.Error(err))
			return err
		}
	}

	w, err := cli.Data()
	if err != nil {
		l.Logger().Error("sendMail send DATA command failed", zap.Error(err))
		return err
	}
	defer w.Close()

	header := make(map[string]string)
	header["From"] = mail.From
	header["To"] = mail.To
	header["Date"] = time.Now().Format("2006-01-02 15:04:05")
	header["Subject"] = mail.Title

	var msg []byte
	for k, v := range header {
		msg = append(msg, []byte(fmt.Sprintf("%s:%s\r\n", k, v))...)
	}
	msg = append(msg, []byte("\r\n")...)
	msg = append(msg, mail.Data...)
	_, err = w.Write(msg)
	if err != nil {
		l.Logger().Error("sendMail write to network failed", zap.Error(err))
		return err
	}

	cli.Quit()
	return nil
}

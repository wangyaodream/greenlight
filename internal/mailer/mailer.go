package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)


//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
    dialer *mail.Dialer
    sender string
}

func New(host string, port int, username, password, sender string) Mailer {
    dialer := mail.NewDialer(host, port, username, password)
    dialer.Timeout = 5 * time.Second

    return Mailer{
        dialer: dialer,
        sender: sender,
    }
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
    tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
    if err != nil {
        return err
    }

    subject := new(bytes.Buffer)
    // 渲染模板并将结果写入到subject缓冲区
    err = tmpl.ExecuteTemplate(subject, "subject", data)
    if err != nil {
        return err
    }

    // 渲染模板并将结果写入到plainBody缓冲区
    plainBody := new(bytes.Buffer)
    err = tmpl.ExecuteTemplate(plainBody, "htmlBody", data)
    if err != nil {
        return err
    }


    htmlBody := new(bytes.Buffer)
    err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
    if err != nil {
        return err
    }

    msg := mail.NewMessage()
    msg.SetHeader("To", recipient)
    msg.SetHeader("From", m.sender)
    msg.SetHeader("Subject", subject.String())
    msg.SetBody("text/plain", plainBody.String())
    msg.AddAlternative("text/html", htmlBody.String())

    // 发送邮件到SMTP服务器
    err = m.dialer.DialAndSend(msg)
    if err != nil {
        return err
    }

    return nil
}

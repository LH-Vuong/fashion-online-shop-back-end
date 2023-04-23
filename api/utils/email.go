package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/initializers"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type VerifyEmailData struct {
	Token    string
	Fullname string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendVerifyEmail(user *model.User, data *VerifyEmailData, templateDir string) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir(templateDir)
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.Execute(&body, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Account Verification")
	m.SetBody("text/html", body.String())
	//m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}

}

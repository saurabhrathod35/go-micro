package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"time"

	"text/template"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}
	data := map[string]any{
		"message": msg.Data,
	}
	msg.Data = data
	log.Println("Sending email to", msg.To, "from", msg.From, "with subject", msg.Subject, msg.Data)
	formattedMessage, err := m.buildHtmlMessage(msg)
	if err != nil {
		log.Println("Error building HTML message", err)
	}
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		log.Println("Error building plain text message", err)
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println("Error connecting to SMTP server", err)
		return err
	}
	email := mail.NewMSG()
	email.SetFrom(msg.FromName+" <"+msg.From+">").
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}
	err = email.Send(smtpClient)
	if err != nil {
		log.Println("Error while sending email ", err)
		return err
	}
	return nil
}

func (m *Mail) buildHtmlMessage(msg Message) (string, error) {
	// Implement the logic to build the HTML message
	// This is a placeholder for the actual implementation
	templateToRender := "./templates/mail.html.gohtml"
	t, err := template.New("body").ParseFiles(templateToRender)
	if err != nil {
		log.Println("Error parsing template:", err)
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Println("Error executing template:", err)
		return "", err
	}
	formattedMessage := tpl.String()
	formattedMessage, err = m.InlineCss(formattedMessage)
	if err != nil {
		log.Println("Error inlining CSS:", err)
		return "", err
	}
	return formattedMessage, nil
}

func (m *Mail) InlineCss(s string) (string, error) {
	// Implement the logic to inline CSS
	// This is a placeholder for the actual implementation
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		log.Println("Error creating premailer:", err)
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		log.Println("Error transforming HTML:", err)
		return "", err
	}
	return html, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// Implement the logic to build the HTML message
	// This is a placeholder for the actual implementation
	templateToRender := "./templates/mail.plain.html.gohtml"
	t, err := template.New("body").ParseFiles(templateToRender)
	if err != nil {
		log.Println("Error parsing template:", err)
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Println("Error executing template:", err)
		return "", err
	}
	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "TLS", "tls":
		return mail.EncryptionSTARTTLS
	case "SSL", "ssl":
		return mail.EncryptionSSLTLS
	case "None", "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Port:        port,
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}
	return m
}

package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func SendEmail(to []string, subject, server, erro, data, path string) {
	from := "prcn29@gmail.com"
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		panic("Erro na variavel de ambiente")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles(path)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	t.Execute(&body, struct {
		Server  string
		Error   string
		Horario string
	}{
		Server:  server,
		Error:   erro,
		Horario: data,
	})

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println("Erro ao enviar email")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Email enviado com sucesso")
}

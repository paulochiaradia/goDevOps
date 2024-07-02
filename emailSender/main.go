package main

import "github.com/paulochiaradia/goDevOps/emailSender/email"

func main() {
	email.SendEmail([]string{"paulo.chiaradia@outlook.com"}, "Alerta: Server down", "Google",
		"Erro ao conectar servidor", "28/06/2024 21:21", "./email/template.html")
}

package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/user/golang-api-rest/internal/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendVerificationEmail(to, name, token string) {
	subject := "Verifique seu email"
	body := fmt.Sprintf(
		"Olá %s,\n\nObrigado por se registrar!\n\n"+
			"Use o token abaixo para verificar seu email:\n%s\n\n"+
			"Ou acesse o link: http://localhost:8080/api/v1/auth/verify-email?token=%s\n\n"+
			"Atenciosamente,\nEquipe API",
		name, token, token,
	)
	s.send(to, subject, body)
}

func (s *EmailService) SendWelcomeEmail(to, name string) {
	subject := "Bem-vindo!"
	body := fmt.Sprintf(
		"Olá %s,\n\nSua conta foi verificada com sucesso!\n\nAtenciosamente,\nEquipe API",
		name,
	)
	s.send(to, subject, body)
}

func (s *EmailService) send(to, subject, body string) {
	if s.cfg.IsDevelopment() {
		log.Printf("[DEV EMAIL] To: %s | Subject: %s\n%s", to, subject, body)
		return
	}

	from := s.cfg.SMTPUser
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		log.Printf("Erro ao enviar email para %s: %v", to, err)
	}
}
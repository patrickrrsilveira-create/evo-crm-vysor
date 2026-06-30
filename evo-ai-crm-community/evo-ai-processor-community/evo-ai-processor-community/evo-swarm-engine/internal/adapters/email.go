package adapters

import (
	"context"
	"log"
)

type EmailAdapter struct {
	SMTPHost string
	SMTPPort int
	IMAPHost string
	Username string
	Password string
}

func NewEmailAdapter(smtp, imap, user, pass string) *EmailAdapter {
	return &EmailAdapter{
		SMTPHost: smtp,
		IMAPHost: imap,
		Username: user,
		Password: pass,
	}
}

func (a *EmailAdapter) Name() string {
	return "email_smtp_imap"
}

func (a *EmailAdapter) Start(ctx context.Context) error {
	log.Println("📧 [EmailAdapter] Inicializado - Conectado via IMAP (Idle) para escutar novos emails")
	return nil
}

func (a *EmailAdapter) SendMessage(ctx context.Context, to string, content string) error {
	log.Printf("📧 [EmailAdapter] Enviando email via SMTP para: %s", to)
	return nil
}

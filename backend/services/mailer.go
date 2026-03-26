package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"upcycleconnect/backend/config"
)

type Mailer struct {
	cfg *config.Config
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: UpcycleConnect <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nDate: %s\r\n\r\n%s",
		m.cfg.MailFrom,
		to,
		subject,
		time.Now().Format(time.RFC1123Z),
		body,
	)

	addr := net.JoinHostPort(m.cfg.SMTPHost, m.cfg.SMTPPort)
	auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPassword, m.cfg.SMTPHost)

	if m.cfg.SMTPPort == "465" {
		return m.sendSSL(addr, auth, to, []byte(msg))
	}
	return smtp.SendMail(addr, auth, m.cfg.MailFrom, []string{to}, []byte(msg))
}

// sendSSL connecte via SSL/TLS direct (port 465)
func (m *Mailer) sendSSL(addr string, auth smtp.Auth, to string, msg []byte) error {
	tlsConfig := &tls.Config{ServerName: m.cfg.SMTPHost}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls.Dial: %w", err)
	}
	client, err := smtp.NewClient(conn, m.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp.NewClient: %w", err)
	}
	defer client.Close()
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth: %w", err)
	}
	if err = client.Mail(m.cfg.MailFrom); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return w.Close()
}

// ─── Templates ────────────────────────────────────────────────────────────────

func (m *Mailer) SendRegisterVerification(to, firstName, token string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", m.cfg.AppURL, token)
	body := fmt.Sprintf(`Bonjour %s,

Merci de vous être inscrit(e) sur UpcycleConnect.

Pour activer votre compte, cliquez sur le lien ci-dessous :

  %s

Ce lien est valable 24 heures.

Si vous n'êtes pas à l'origine de cette inscription, vous pouvez ignorer cet email.

— L'équipe UpcycleConnect`, firstName, link)

	return m.Send(to, "Activez votre compte UpcycleConnect", body)
}

func (m *Mailer) SendEmailChangeCode(to, firstName, code string) error {
	body := fmt.Sprintf(`Bonjour %s,

Vous avez demandé à modifier l'adresse email de votre compte UpcycleConnect.

Votre code de vérification est :

  %s

Ce code est valable 10 minutes.

Si vous n'êtes pas à l'origine de cette demande, ignorez cet email et votre adresse ne sera pas modifiée.

— L'équipe UpcycleConnect`, firstName, code)

	return m.Send(to, "Code de vérification — Changement d'email", body)
}

func (m *Mailer) SendLoginVerification(to, firstName, token, ip, attemptTime string) error {
	link := fmt.Sprintf("%s/verify-login?token=%s", m.cfg.AppURL, token)
	body := fmt.Sprintf(`Bonjour %s,

Une tentative de connexion à votre compte a été détectée depuis un emplacement inconnu.

  Date et heure : %s
  Adresse IP    : %s

Pour confirmer cette connexion, cliquez sur le lien ci-dessous :

  %s

Ce lien est valable 15 minutes.

Si vous n'êtes pas à l'origine de cette tentative, ignorez cet email et changez votre mot de passe immédiatement.

— L'équipe UpcycleConnect`, firstName, attemptTime, ip, link)

	return m.Send(to, "Nouvelle connexion — Confirmez votre identité", body)
}

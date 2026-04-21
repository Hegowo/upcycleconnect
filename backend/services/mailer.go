package services

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
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

func (m *Mailer) SendWithAttachment(to, subject, body, attachmentPath string) error {
	if attachmentPath == "" {
		return m.Send(to, subject, body)
	}

	data, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("read attachment: %w", err)
	}

	boundaryBytes := make([]byte, 16)
	if _, err := rand.Read(boundaryBytes); err != nil {
		return fmt.Errorf("boundary: %w", err)
	}
	boundary := fmt.Sprintf("uc-%x", boundaryBytes)
	filename := filepath.Base(attachmentPath)

	var b strings.Builder
	fmt.Fprintf(&b, "From: UpcycleConnect <%s>\r\n", m.cfg.MailFrom)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(&b, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary)

	fmt.Fprintf(&b, "--%s\r\n", boundary)
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	b.WriteString(body)
	b.WriteString("\r\n\r\n")

	fmt.Fprintf(&b, "--%s\r\n", boundary)
	fmt.Fprintf(&b, "Content-Type: application/pdf; name=\"%s\"\r\n", filename)
	b.WriteString("Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(&b, "Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename)

	encoded := base64.StdEncoding.EncodeToString(data)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		b.WriteString(encoded[i:end])
		b.WriteString("\r\n")
	}
	fmt.Fprintf(&b, "--%s--\r\n", boundary)

	addr := net.JoinHostPort(m.cfg.SMTPHost, m.cfg.SMTPPort)
	auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPassword, m.cfg.SMTPHost)
	msg := []byte(b.String())

	if m.cfg.SMTPPort == "465" {
		return m.sendSSL(addr, auth, to, msg)
	}
	return smtp.SendMail(addr, auth, m.cfg.MailFrom, []string{to}, msg)
}

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

func (m *Mailer) SendPaymentConfirmation(to, firstName, prestationTitle, invoiceNumber string, amountCents int64, pdfPath string) error {
	amount := fmt.Sprintf("%.2f EUR", float64(amountCents)/100.0)
	body := fmt.Sprintf(`Bonjour %s,

Nous vous confirmons la réception de votre paiement pour la prestation :

  %s

Montant réglé : %s
Numéro de facture : %s

Vous trouverez votre facture en pièce jointe de cet email. Elle est également disponible à tout moment depuis votre espace personnel, rubrique "Mes factures".

Merci de votre confiance.

— L'équipe UpcycleConnect`, firstName, prestationTitle, amount, invoiceNumber)

	return m.SendWithAttachment(to, "Confirmation de paiement — UpcycleConnect", body, pdfPath)
}

func (m *Mailer) SendQuote(to, firstName, prestationTitle, quoteNumber, pdfPath string) error {
	body := fmt.Sprintf(`Bonjour %s,

Suite à votre demande concernant la prestation :

  %s

Vous trouverez ci-joint votre devis n° %s.

Ce devis est valable 30 jours à compter de sa date d'émission. Il est également disponible dans votre espace personnel, rubrique "Mes factures".

Pour toute question, n'hésitez pas à nous contacter.

— L'équipe UpcycleConnect`, firstName, prestationTitle, quoteNumber)

	return m.SendWithAttachment(to, "Votre devis UpcycleConnect", body, pdfPath)
}

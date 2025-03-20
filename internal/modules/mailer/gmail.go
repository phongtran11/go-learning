package mailer

import (
	"fmt"
	"modular-fx-fiber/internal/core/config"

	"github.com/wneessen/go-mail"
)

// GmailMailer is a mailer configured for Gmail with template support
type GmailMailer struct {
	client     *mail.Client
	from       string
	templates  *TemplateManager
	defaultCtx map[string]interface{}
}

// GmailMailerConfig holds configuration for the Gmail mailer
type GmailMailerConfig struct {
	GmailAddress string
	Password     string
	FromName     string
	TemplateDir  string
}

// NewGmailMailer creates a new Gmail-configured mailer with template support
func NewGmailMailer(c *config.Config, tm *TemplateManager) (*GmailMailer, error) {
	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(c.Mail.SMTPUsername),
		mail.WithPassword(c.Mail.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory))

	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	mailer := &GmailMailer{
		client:     client,
		from:       c.Mail.FromAddr,
		templates:  tm,
		defaultCtx: make(map[string]any),
	}

	return mailer, nil
}

// SetDefaultContext sets default context values for all templates
func (g *GmailMailer) SetDefaultContext(key string, value any) {
	g.defaultCtx[key] = value
}

// SendEmail sends a basic email through Gmail
func (g *GmailMailer) SendEmail(to, subject, textBody, htmlBody string) error {
	msg := mail.NewMsg()

	if err := msg.From(g.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(subject)

	// Set the body - HTML with plain text alternative if both are provided
	if htmlBody != "" {
		msg.SetBodyString(mail.TypeTextHTML, htmlBody)
		if textBody != "" {
			msg.AddAlternativeString(mail.TypeTextPlain, textBody)
		}
	} else if textBody != "" {
		msg.SetBodyString(mail.TypeTextPlain, textBody)
	}

	// Send the email
	return g.client.DialAndSend(msg)
}

// SendTemplatedEmail sends an email using a template
func (g *GmailMailer) SendTemplatedEmail(to, subject, templateName string, ctx map[string]any) error {
	if g.templates == nil {
		return fmt.Errorf("template manager not initialized")
	}

	// Merge default context with provided context
	mergedCtx := make(map[string]any)
	for k, v := range g.defaultCtx {
		mergedCtx[k] = v
	}
	for k, v := range ctx {
		mergedCtx[k] = v
	}

	// Render the template
	htmlContent, err := g.templates.RenderTemplate(templateName, mergedCtx)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send the email
	return g.SendEmail(to, subject, "", htmlContent)
}

// Close closes the mail client
func (g *GmailMailer) Close() error {
	return g.client.Close()
}

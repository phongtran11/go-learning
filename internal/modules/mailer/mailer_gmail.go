package mailer

import (
	"fmt"
	"maps"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/shared/logger"

	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

// GmailMailer is a mailer configured for Gmail with template support
type (
	GmailMailer interface {
		SetDefaultContext(key string, value any)
		SendEmail(to, subject, textBody, htmlBody string) error
		SendTemplatedEmail(to, subject, templateName string, ctx map[string]any) error
		Close() error
	}

	gmailMailer struct {
		client     *mail.Client
		from       string
		templates  *TemplateManager
		logger     *logger.ZapLogger
		defaultCtx map[string]any
	}
)

// NewGmailMailer creates a new Gmail-configured mailer with template support
func NewGmailMailer(l *logger.ZapLogger, c *config.Config, tm *TemplateManager) GmailMailer {
	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(c.Mail.SMTPUsername),
		mail.WithPassword(c.Mail.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory))

	if err != nil {
		l.Error("failed to create mail client", zap.Error(err))
		return nil
	}

	return &gmailMailer{
		logger:     l,
		client:     client,
		from:       c.Mail.FromAddr,
		templates:  tm,
		defaultCtx: make(map[string]any),
	}

}

// SetDefaultContext sets default context values for all templates
func (g *gmailMailer) SetDefaultContext(key string, value any) {
	g.defaultCtx[key] = value
}

// SendEmail sends a basic email through Gmail
func (g *gmailMailer) SendEmail(to, subject, textBody, htmlBody string) error {
	g.logger.Debug("Preparing to send email",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.Int("textBodyLength", len(textBody)),
		zap.Int("htmlBodyLength", len(htmlBody)))

	msg := mail.NewMsg()

	if err := msg.From(g.from); err != nil {
		g.logger.Error("Failed to set from address",
			zap.String("from", g.from),
			zap.Error(err))
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(to); err != nil {
		g.logger.Error("Failed to set to address",
			zap.String("to", to),
			zap.Error(err))
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

	// Log before sending
	g.logger.Info("Attempting to send email",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("from", g.from))

	// Send the email
	err := g.client.DialAndSend(msg)
	if err != nil {
		g.logger.Error("Failed to send email",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}

	g.logger.Info("Email sent successfully",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	return nil
}

// SendTemplatedEmail sends an email using a template
func (g *gmailMailer) SendTemplatedEmail(to, subject, templateName string, ctx map[string]any) error {

	if g.templates == nil {
		return fmt.Errorf("template manager not initialized")
	}
	sendMailData := fmt.Sprintf("email %v , subject %v, template %v", to, subject, templateName)
	g.logger.Info("Send Mail Data:", zap.Any("sendMailData", sendMailData))

	// Merge default context with provided context
	mergedCtx := make(map[string]any)
	maps.Copy(mergedCtx, g.defaultCtx)
	maps.Copy(mergedCtx, ctx)

	// Render the template
	htmlContent, err := g.templates.RenderTemplate(templateName, mergedCtx)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send the email
	return g.SendEmail(to, subject, "", htmlContent)
}

// Close closes the mail client
func (g *gmailMailer) Close() error {
	return g.client.Close()
}

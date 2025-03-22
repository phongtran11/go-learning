package mailer

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTemplateManager),
	fx.Provide(NewGmailMailer),
)

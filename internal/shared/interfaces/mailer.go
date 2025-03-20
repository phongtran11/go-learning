package interfaces

type MailerService interface {
	SendHtmlMail(to string, subject string, htmlBody string) error
}

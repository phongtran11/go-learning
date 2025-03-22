package mailer

const (
	// EmailVerificationSubject is the subject of the email verification email
	EmailVerificationSubject  = "Email Verification"
	EmailVerificationTemplate = "send_confirm_email_code"
)

type EmailVerificationData struct {
	Name string
	Code string
}

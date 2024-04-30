package messaging

type (
	Messaging interface {
		Push(Mail) error
	}

	Author struct {
		name  string
		email string
	}

	Mail struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Subject  string `json:"subject"`
		TextPart string `json:"textPart"`
		HtmlPart string `json:"htmlPart"`
	}
)

func NewMailer(name string, email string) *Author {
	return &Author{
		name:  name,
		email: email,
	}
}

func BuildMail(name, email, subject, template string) Mail {
	return Mail{
		Name:     name,
		Email:    email,
		Subject:  subject,
		HtmlPart: template,
	}
}

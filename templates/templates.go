package templates

import (
	"errors"
	"fmt"
)

const (
	FORGOT_PASSWORD = "FORGOT_PASSWORD"

	ACCOUNT_CREATED = "ACCOUNT_CREATED"
)

func NewTemplate(templateKey string, args ...any) (string, error) {
	switch templateKey {
	case FORGOT_PASSWORD:
		return fmt.Sprintf(ForgortPasswordTemplate, args...), nil
	case ACCOUNT_CREATED:
		return fmt.Sprintf(AccountCreatedTemplate, args...), nil
	default:
		return "", errors.New("invalid template key")
	}
}

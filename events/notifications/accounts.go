package notifications

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/tejiriaustin/narx_api/consumer"
	"github.com/tejiriaustin/narx_api/events"
	"github.com/tejiriaustin/narx_api/messaging"
	"github.com/tejiriaustin/narx_api/templates"
)

const (
	ForgotPasswordNotification = "NOTIFICATION.FORGOT_PASSWORD"
)

func ForgotPasswordNotificationEventHandler(mailer messaging.Messaging) consumer.Handler {
	return func(ctx context.Context, msg events.Event) error {
		fmt.Println("Forgot Password")

		template, err := templates.NewTemplate(templates.FORGOT_PASSWORD, msg.MsgBody["full_name"], msg.MsgBody["code"])
		if err != nil {
			zap.L().Error("failed to create template for forgot password", zap.String("template", template), zap.Any("data", msg))
			return errors.New("failed to send forgot password email")
		}
		email := msg.MsgBody["email"].(string)

		err = mailer.Push(email, template)
		if err != nil {
			zap.L().Error("failed to push mail", zap.Error(err))
			return err
		}

		return nil
	}
}

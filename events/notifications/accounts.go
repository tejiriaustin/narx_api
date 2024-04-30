package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tejiriaustin/slabmark-api/events"

	"go.uber.org/zap"

	"github.com/tejiriaustin/slabmark-api/consumer"
	"github.com/tejiriaustin/slabmark-api/messaging"
	"github.com/tejiriaustin/slabmark-api/templates"
)

const (
	ForgotPasswordNotification = "NOTIFICATION.FORGOT_PASSWORD"
)

func ForgotPasswordNotificationEventHandler(mailer messaging.Messaging) consumer.Handler {
	return func(ctx context.Context, msg consumer.Message) error {
		fmt.Println("Forgot Password")

		var data events.ForgotPasswordNotificationSchema
		err := json.Unmarshal([]byte(msg.Body), &data)
		if err != nil {
			zap.L().Error("failed to unmarshal message", zap.Error(err))
			return err
		}

		template, err := templates.NewTemplate(templates.FORGOT_PASSWORD)
		if err != nil {
			zap.L().Error("failed to create template for forgot password", zap.String("template", template), zap.Any("data", data))
			return errors.New("failed to send forgot password email")
		}

		mail := messaging.BuildMail(
			data.FirstName,
			data.Email,
			"Forgot your password",
			template,
		)
		err = mailer.Push(mail)
		if err != nil {
			zap.L().Error("failed to push mail", zap.Error(err))
			return err
		}

		return nil
	}
}

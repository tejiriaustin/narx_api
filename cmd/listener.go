package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/tejiriaustin/narx_api/consumer"
	"github.com/tejiriaustin/narx_api/database"
	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/events/notifications"
	"github.com/tejiriaustin/narx_api/messaging"
)

// apiCmd represents the api command
var listenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Starts narx-api consumer service",
	Long:  ``,
	Run:   startListener,
}

func init() {
	rootCmd.AddCommand(listenerCmd)
}

func startListener(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config := setListenerEnvironment()

	dbConn, err := database.NewMongoDbClient().Connect(config.GetAsString(env.MongoDsn), config.GetAsString(env.MongoDbName))
	if err != nil {
		panic("Couldn't connect to mongo dsn: " + err.Error())
	}
	defer func() {
		_ = dbConn.Disconnect(context.TODO())
	}()

	mailer := messaging.NewSMTP(
		config.GetAsString(env.SmtpPassword),
		config.GetAsString(env.SmtpSender),
		config.GetAsString(env.SmtpHost),
		config.GetAsString(env.SmtpPort),
	)
	db := dbConn.GetCollection("notifications")

	listeners := consumer.NewConsumer(consumer.WithUpdater(db)).
		SetHandler(notifications.ForgotPasswordNotification, notifications.ForgotPasswordNotificationEventHandler(mailer))

	listeners.ListenAndServe(ctx, db)
}

func setListenerEnvironment() env.Environment {
	staticEnvironment := env.NewEnvironment()

	staticEnvironment.
		SetEnv(env.MongoDsn, env.MustGetEnv(env.MongoDsn)).
		SetEnv(env.MongoDbName, env.MustGetEnv(env.MongoDbName)).
		SetEnv(env.SmtpHost, env.MustGetEnv(env.SmtpHost)).
		SetEnv(env.SmtpSender, env.MustGetEnv(env.SmtpSender)).
		SetEnv(env.SmtpPort, env.MustGetEnv(env.SmtpPort)).
		SetEnv(env.SmtpPassword, env.MustGetEnv(env.SmtpPassword)).
		SetEnv(env.FirebaseAuthKey, env.MustGetEnv(env.FirebaseAuthKey)).
		SetEnv(env.FirebaseRegistrationToken, env.MustGetEnv(env.FirebaseRegistrationToken)).
		SetEnv(env.FirebaseServiceAccountKey, env.MustGetEnv(env.FirebaseServiceAccountKey))

	return staticEnvironment
}

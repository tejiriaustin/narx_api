package cmd

import (
	"context"
	"github.com/tejiriaustin/narx_api/publisher"

	"github.com/spf13/cobra"

	"github.com/tejiriaustin/narx_api/database"
	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/server"
	"github.com/tejiriaustin/narx_api/services"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts NARX PV api",
	Long:  ``,
	Run:   startApi,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func startApi(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config := setApiEnvironment()

	dbConn, err := database.NewMongoDbClient().Connect(config.GetAsString(env.MongoDsn), config.GetAsString(env.MongoDbName))
	if err != nil {
		panic("Couldn't connect to mongo dsn: " + err.Error())
	}
	defer func() {
		_ = dbConn.Disconnect(context.TODO())
	}()

	rc := repository.NewRepositoryContainer(dbConn)

	sc := services.NewService(&config)
	sc.Publisher = publisher.NewPublisher(dbConn.GetCollection("notifications"))

	server.Start(ctx, sc, rc, &config)
}

func setApiEnvironment() env.Environment {
	staticEnvironment := env.NewEnvironment()

	staticEnvironment.
		SetEnv(env.Port, env.GetEnv(env.Port, "8080")).
		SetEnv(env.MongoDsn, env.MustGetEnv(env.MongoDsn)).
		SetEnv(env.MongoDbName, env.MustGetEnv(env.MongoDbName)).
		SetEnv(env.JwtSecret, env.MustGetEnv(env.JwtSecret)).
		SetEnv(env.FrontendUrl, env.MustGetEnv(env.FrontendUrl))

	return staticEnvironment
}

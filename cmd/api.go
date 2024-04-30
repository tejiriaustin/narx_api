package cmd

import (
	"context"

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

	redis, err := database.NewRedisClient(config.GetAsString(env.RedisDsn), config.GetAsString(env.RedisPass), config.GetAsString(env.RedisPort))
	if err != nil {
		panic("Couldn't connect to redis dsn: " + err.Error())
	}
	defer func() {
		_ = redis.Disconnect(context.TODO())
	}()

	rc := repository.NewRepositoryContainer(dbConn)

	sc := services.NewService(&config)

	server.Start(ctx, sc, rc, &config)
}

func setApiEnvironment() env.Environment {
	staticEnvironment := env.NewEnvironment()

	staticEnvironment.
		SetEnv(env.EnvPort, env.GetEnv(env.EnvPort, "8080")).
		SetEnv(env.RedisDsn, env.MustGetEnv(env.RedisDsn)).
		SetEnv(env.RedisPort, env.MustGetEnv(env.RedisPort)).
		SetEnv(env.RedisPass, env.MustGetEnv(env.RedisPass)).
		SetEnv(env.MongoDsn, env.MustGetEnv(env.MongoDsn)).
		SetEnv(env.MongoDbName, env.MustGetEnv(env.MongoDbName)).
		SetEnv(env.JwtSecret, env.MustGetEnv(env.JwtSecret))

	return staticEnvironment
}

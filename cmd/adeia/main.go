package main

import (
	"fmt"
	"os"

	"adeia/internal/cache/redis"
	"adeia/internal/config"
	"adeia/internal/http"
	"adeia/internal/http/server"
	"adeia/internal/repo"
	"adeia/internal/service"
	"adeia/internal/store/pg"
	"adeia/pkg/log/zap"
	"adeia/pkg/util/constants"
	"adeia/pkg/util/ioutil"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgres driver
)

func main() {
	confPath := getEnv(constants.EnvConfPathKey, "config/config.yaml")
	confFile, err := os.Open(confPath)
	checkErr(err)

	conf, err := config.Load(confFile)
	checkErr(err)

	checkErr(run(conf))
}

func run(conf *config.Config) error {
	logger, err := zap.New(&conf.LoggerConfig)
	if err != nil {
		return err
	}

	dbConn, cacheConn, err := initConnections(conf)
	if err != nil {
		logger.Debugf("failed to initialize connections: %v", err)
		return err
	}

	defer func() {
		logger.Debug("closing connections...")
		_ = logger.Sync()
		ioutil.CheckCloseErr(dbConn, &err)
		ioutil.CheckCloseErr(cacheConn, &err)
	}()

	// init repos
	logger.Debug("initializing repositories...")
	userRepo := repo.NewUserRepo(dbConn)

	// init services
	logger.Debug("initializing services...")
	userService := service.NewUserService(logger, userRepo)

	// init controllers
	logger.Debug("initializing controllers...")
	userController := http.NewUserController(logger, userService)

	srv := server.New(&conf.ServerConfig, logger, userController)
	srv.BindControllers()
	srv.Serve()

	return nil
}

func initConnections(conf *config.Config) (*pg.PostgresDB, *redis.Redis, error) {
	dbConn, err := pg.New(&conf.DBConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize connection to db: %v", err)
	}

	cacheConn, err := redis.New(&conf.CacheConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize connection to cache: %v", err)
	}

	return dbConn, cacheConn, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

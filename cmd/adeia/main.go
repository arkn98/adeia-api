/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

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
	"adeia/pkg/constants"
	"adeia/pkg/log/zap"
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
	log, err := zap.New(&conf.LoggerConfig)
	if err != nil {
		return err
	}

	db, cache, err := initConnections(&conf.DBConfig, &conf.CacheConfig)
	if err != nil {
		log.Errorf("failed to initialize connections: %v", err)
		return err
	}

	defer func() {
		log.Info("closing connections...")
		err = log.Sync()
		ioutil.CheckCloseErr(db, &err)
		ioutil.CheckCloseErr(cache, &err)
	}()

	// init repos
	log.Debug("initializing repositories...")
	userRepo := repo.NewUserRepo(db)

	// init services
	log.Debug("initializing services...")
	userService := service.NewUserService(log, userRepo)

	// init controllers
	log.Debug("initializing controllers...")
	userController := http.NewUserController("/users", log, userService)

	srv := server.New(&conf.ServerConfig, log, userController)
	srv.Start()

	return err
}

func initConnections(d *config.DBConfig, c *config.CacheConfig) (*pg.PostgresDB, *redis.Redis, error) {
	db, err := pg.New(d)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize db connection: %v", err)
	}

	cache, err := redis.New(c)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot initialize cache connection: %v", err)
	}

	return db, cache, nil
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

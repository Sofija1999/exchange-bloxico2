// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"fmt"
	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/Bloxico/exchange-gateway/sofija/log"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"os"
)

// Injectors from wire.go:

// Injectors
func InitializeApp() (*App, error) {
	fileProvider := config.DefaultFileProvider()
	configConfig, err := config.NewConfig(fileProvider)
	if err != nil {
		return nil, err
	}
	environment := envConfigProvider(configConfig)
	logger := initializeLogger(environment)
	db, err := initializeDatabase(configConfig)
	if err != nil {
		return nil, err
	}
	app := &App{
		Config: configConfig,
		Logger: logger,
		DB:     db,
	}
	return app, nil
}

func InitializeTestApp() (*App, error) {
	fileProvider := config.TestFileProvider()
	configConfig, err := config.NewConfig(fileProvider)
	if err != nil {
		return nil, err
	}
	environment := envConfigProvider(configConfig)
	logger := initializeLogger(environment)
	db, err := initializeDatabase(configConfig)
	if err != nil {
		return nil, err
	}
	app := &App{
		Config: configConfig,
		Logger: logger,
		DB:     db,
	}
	return app, nil
}

// wire.go:

// Providers
func envConfigProvider(cfg config.Config) config.Environment { return cfg.Env }

func httpConfigProvider(cfg config.Config) config.ServerConfig { return cfg.Http }

func dbConfigProvider(cfg config.Config) config.DatabaseConfig { return cfg.Database }

var ConfigSet = wire.NewSet(config.FileProviderSet, config.NewConfig)

var TestConfigSet = wire.NewSet(config.TestFileProviderSet, config.NewConfig)

var LoggerSet = wire.NewSet(
	envConfigProvider,
	initializeLogger,
)

var DatabaseSet = wire.NewSet(
	dbConfigProvider,
	initializeDatabase,
)

var EgwUserSet = wire.NewSet(repo.NewEgwUserRepository)

// All config vars for the App
var AppSet = wire.NewSet(config.NewConfig, LoggerSet,
	DatabaseSet, wire.Struct(new(App), "*"),
)

func MustInitializeApp() *App {
	app, err := InitializeApp()
	if err != nil {
		fmt.Printf("failed initializing app: %s", err)
		os.Exit(1)
		return nil
	}

	return app
}

func MustInitializeTestApp() *App {
	app, err := InitializeTestApp()
	if err != nil {
		fmt.Printf("failed initializing test app: %s", err)
		os.Exit(1)
		return nil
	}

	return app
}

func initializeLogger(env config.Environment) log.Logger {
	logger, err := log.NewLogger(env == config.EnvLocal)
	if err != nil {
		fmt.Println(errors.Wrap(err, "config logger"))
		os.Exit(1)
	}

	return logger
}

func initializeDatabase(cfg config.Config) (*database.DB, error) {
	return database.NewDB(cfg.Database)
}

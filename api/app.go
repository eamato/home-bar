package api

import (
	"context"
	"home-bar/configs"
	"home-bar/database"
	"home-bar/internal"
	"os"
	"os/signal"
	"syscall"
)

// Application struct holds all required things to work, config, DI database, etc.
type Application struct {
	Config         *configs.Config
	DatabaseClient database.Client
}

func App() Application {
	cfg := configs.NewConfig()
	internal.RegisterValidators()
	return Application{
		Config:         cfg,
		DatabaseClient: database.NewHomeBarDatabase(cfg),
	}
}

func (app *Application) CloseApp(onAppCloses func()) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGTERM, syscall.SIGINT)
	<-quitChannel

	if app.DatabaseClient != nil {
		err := app.DatabaseClient.Disconnect(context.TODO())
		if err != nil {
			internal.PrintFatal("", err)
		}

		internal.PrintMessage("Connection to database closed.")
	}

	onAppCloses()
}

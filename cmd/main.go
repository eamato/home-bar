package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"home-bar/api"
	"home-bar/api/route"
	"home-bar/internal"
	"time"
)

//https://github.com/go-sql-driver/mysql

func main() {
	app := api.App()

	ginEngine := gin.Default()
	ginEngine.LoadHTMLGlob("static/templates/*.html")

	apiGroup := ginEngine.Group("api")
	webGroup := ginEngine.Group("web")

	database := app.DatabaseClient.Database(app.Config.DatabaseConfig.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.DatabaseClient.Ping(ctx); err != nil {
		internal.PrintFatal("error occurred while pinging database", err)
	}

	route.Setup(app.Config, database, apiGroup)
	route.SetupWeb(app.Config, database, webGroup)

	server := api.NewServer(app.Config, ginEngine.Handler())

	go func() {
		if err := server.RunServer(); err != nil {
			internal.PrintFatal("error occurred while running http server", err)
		}
	}()

	app.CloseApp(func() {
		if err := server.ShutDown(context.TODO()); err != nil {
			internal.PrintFatal("error occurred while server shutdown", err)
		}
	})
}

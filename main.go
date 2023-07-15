package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ianugroho1994/todo/project"
	"github.com/ianugroho1994/todo/shared"
	"github.com/ianugroho1994/todo/task"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "c", "config.yml", "Config file name")

	flag.Parse()

	cfg := defaultConfig()
	cfg.loadFromEnv()

	if len(configFileName) > 0 {
		err := loadConfigFromFile(configFileName, &cfg)
		if err != nil {
			log.Warn().Str("file", configFileName).Err(err).Msg("cannot load config file, use defaults")
		}
	}

	log.Debug().Any("config", cfg).Msg("config loaded")

	dbConn, err := sqlx.Connect(`mysql`, cfg.DBConfig.ConnStr())
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot ping db")
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	shared.SetDBConnection(dbConn)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Mount("/todo", task.TaskRouters())
	r.Mount("/todo", project.ProjectRouters())

	log.Info().Msg("Starting up server...")

	if err := http.ListenAndServe(cfg.Listen.Addr(), r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return
	}

	log.Info().Msg("Server Stopped")
}

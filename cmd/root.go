package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rudbast/my-net/core"
	"github.com/rudbast/my-net/middleware"
	"github.com/rudbast/my-net/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	DBOption struct {
		Username string
		Password string
		Host     string
		Port     int
		Name     string
	}
)

var (
	router *mux.Router
	db     *sql.DB
	logger *log.Logger

	rootCmd = &cobra.Command{
		Use: "my-net",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Initiate config.
			viper.SetConfigType("toml")

			// Search config in home directory with name "config" (without extension).
			viper.AddConfigPath("./files/data/my-net")
			viper.AddConfigPath("/data/my-net")
			viper.SetConfigName("config")

			// Read in environment variables that match.
			viper.AutomaticEnv()

			logger = log.New(os.Stdout, "", log.LstdFlags)

			// If a config file is found, read it in.
			err := viper.ReadInConfig()
			if err != nil {
				logger.Fatalln("Read config file error:", err)
			}

			// Initiate database.
			db, err := connectDatabase(DBOption{
				Username: viper.GetString("database.username"),
				Password: viper.GetString("database.password"),
				Host:     viper.GetString("database.host"),
				Port:     viper.GetInt("database.port"),
				Name:     viper.GetString("database.name"),
			})
			if err != nil {
				logger.Fatalln("Connect database error:", err)
			}

			service := core.New(logger, db)
			module := handler.New(logger, service)

			// Initiate routes.
			router = mux.NewRouter()

			router.Use(middleware.ContextMiddleware)

			router.Handle("/user/:followed_id/follow", HandlerFunc(module.HandleUserFollow)).Methods(http.MethodPost)
			router.Handle("/user/:followed_id/unfollow", HandlerFunc(module.HandleUserUnfollow)).Methods(http.MethodPost)

			router.Handle("/feeds", HandlerFunc(module.HandleUserFeed)).Methods(http.MethodGet)
			router.Handle("/tweet", HandlerFunc(module.HandleUserPost)).Methods(http.MethodPost)
		},
		Run: func(cmd *cobra.Command, args []string) {
			port := viper.GetInt("app.port")

			srv := &http.Server{
				Addr:    fmt.Sprintf(":%d", port),
				Handler: router,
			}

			idleConnsClosed := make(chan struct{})
			go func() {
				sigint := make(chan os.Signal, 1)
				signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
				<-sigint

				// We received an interrupt signal, shut down.
				if err := srv.Shutdown(context.Background()); err != nil {
					// Error from closing listeners, or context timeout.
					logger.Println("Server shutdown error.")
				}

				logger.Println("Server shutdown.")
				close(idleConnsClosed)
			}()

			logger.Printf("Server listening on :%d.\n", port)
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				logger.Fatalln("Listen and serve error.")
			}

			<-idleConnsClosed
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func connectDatabase(opt DBOption) (*sql.DB, error) {
	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		opt.Username,
		opt.Password,
		opt.Host,
		opt.Port,
		opt.Name)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "db/open")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db/ping")
	}

	return db, nil
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Data  interface{} `json:"data"`
		Error string      `json:"error,omitempty"`
	}{}

	data, err := fn(w, r)
	if err != nil {
		// TODO: Add proper mapping for http status on request failure.
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: Add proper error message (verbose & user-friendly).
		resp.Error = err.Error()
	} else {
		resp.Data = data
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		logger.Println("Encode response error:", err)
	}
}

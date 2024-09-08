package api

import (
	"bank-authentication-system/internal/publisher/handler"
	"bank-authentication-system/pkg/config"
	"bank-authentication-system/pkg/model"
	"bank-authentication-system/pkg/mqtt"
	"bank-authentication-system/pkg/storage/db"
	"bank-authentication-system/pkg/storage/s3"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "start a new api server for bank system",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	})
}

func main(cfg config.Config) {
	database := db.WithRetry(db.Create, cfg.Postgres)
	defer func() {
		if err := database.Close(); err != nil {
			logrus.Error(err.Error())
		}
	}()

	s, err := s3.NewSession(cfg.S3)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	mq, err := mqtt.NewConnection(cfg.MQTT)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	userRepo := model.NewSQLUserRepo(database)

	userHandler := handler.UserHandler{
		UserRepo: userRepo,
		S3:       s,
		MQTT:     mq,
	}

	e := echo.New()
	e.POST("/api", userHandler.RegisterRequestHandler)
	e.GET("/api/:id", userHandler.CheckRequestStatusHandler)

	port := 8080
	serverAddress := fmt.Sprintf(":%d", port)
	if err := e.Start(serverAddress); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Server error: %v", err)
	}
}

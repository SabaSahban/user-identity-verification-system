package consumer

import (
	"bank-authentication-system/internal/consumer/faceAuth"
	"bank-authentication-system/pkg/config"
	"bank-authentication-system/pkg/model"
	"bank-authentication-system/pkg/mqtt"
	"bank-authentication-system/pkg/service/imagga"
	"bank-authentication-system/pkg/service/mail"
	"bank-authentication-system/pkg/storage/db"
	"bank-authentication-system/pkg/storage/s3"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(&cobra.Command{
		Use:   "consumer",
		Short: "starting the face authentication service",
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

	s3, err := s3.NewSession(cfg.S3)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	mqtt, err := mqtt.NewConnection(cfg.MQTT)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	// mailgun connection
	mailGun := mail.NewConnection(cfg.MailGun)

	// imagga handler
	imagga := &imagga.Imagga{Cfg: cfg.Imagga}

	userRepo := model.NewSQLUserRepo(database)

	f := faceAuth.FaceAuth{
		Imagga:   imagga,
		UserRepo: userRepo,
		MailGun:  mailGun,
		S3:       s3,
		MQTT:     mqtt,
	}

	f.Process()
}

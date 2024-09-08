package config

import (
	"time"
)

func Default() Config {
	return Config{
		Postgres: Postgres{
			Host:               "",
			Port:               5432,
			Username:           "",
			Password:           "",
			DBName:             "",
			ConnectTimeout:     30 * time.Second,
			ConnectionLifetime: 30 * time.Minute,
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
			Debug:              true,
		},
		S3: S3{
			AccessKeyID:     "",
			SecretAccessKey: "",
			Region:          "",
			Bucket:          "",
			Endpoint:        "",
		},
		MQTT: MQTT{
			Queue: "",
			URI:   "",
		},
		Imagga: Imagga{
			ApiKey:    "",
			ApiSecret: "",
		},
		MailGun: MailGun{
			Domain: "",
			APIKEY: "",
		},
	}
}

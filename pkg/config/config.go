package config

import (
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/structs"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Postgres Postgres `koanf:"postgres"`
	S3       S3       `koanf:"s3"`
	MQTT     MQTT     `koanf:"mqtt"`
	Imagga   Imagga   `koanf:"imagga"`
	MailGun  MailGun  `koanf:"mail-gun"`
}
type Postgres struct {
	Host               string        `koanf:"host"`
	Port               int           `koanf:"port"`
	Username           string        `koanf:"user"`
	Password           string        `koanf:"pass"`
	DBName             string        `koanf:"dbname"`
	ConnectTimeout     time.Duration `koanf:"connect-timeout"`
	ConnectionLifetime time.Duration `koanf:"connection-lifetime"`
	MaxOpenConnections int           `koanf:"max-open-connections"`
	MaxIdleConnections int           `koanf:"max-idle-connections"`
	Debug              bool          `koanf:"debug"`
}

type S3 struct {
	AccessKeyID     string `koanf:"accessKeyID"`
	SecretAccessKey string `koanf:"secretAccessKey"`
	Region          string `koanf:"region"`
	Bucket          string `koanf:"bucket"`
	Endpoint        string `koanf:"endpoint"`
}

type MQTT struct {
	Queue string `koanf:"queue"`
	URI   string `koanf:"uri"`
}

type Imagga struct {
	ApiKey    string `koanf:"api_key"`
	ApiSecret string `koanf:"api_secret"`
}

type MailGun struct {
	Domain string `koanf:"domain"`
	APIKEY string `koanf:"api_key"`
}

func InitConfig() Config {
	var cfg Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Unmarshal("postgres", &cfg.Postgres); err != nil {
		logrus.Fatalf("error unmarshaling postgres configuration: %s", err)
	}

	if err := k.Unmarshal("mqtt", &cfg.MQTT); err != nil {
		logrus.Fatalf("error unmarshaling postgres configuration: %s", err)
	}

	if err := k.Unmarshal("s3", &cfg.S3); err != nil {
		logrus.Fatalf("error unmarshaling s3 configuration: %s", err)
	}

	if err := k.Unmarshal("mail-gun", &cfg.MailGun); err != nil {
		logrus.Fatalf("error unmarshaling mailgun configuration: %s", err)
	}

	if err := k.Unmarshal("imagga", &cfg.Imagga); err != nil {
		logrus.Fatalf("error unmarshaling imagga configuration: %s", err)
	}

	return cfg
}

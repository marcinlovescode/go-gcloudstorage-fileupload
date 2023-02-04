package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App           `yaml:"app"`
		HTTP          `yaml:"http"`
		Log           `yaml:"logger"`
		GCloudStorage `yaml:"gcloud_storage"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port            string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ReadTimeout     int64  `env-required:"true" yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout    int64  `env-required:"true" yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
		ShutdownTimeout int64  `env-required:"true" yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	GCloudStorage struct {
		ProjectName       string `env-required:"true" yaml:"project_name"    env:"GCLOUD_PROJECT_NAME"`
		BucketName        string `env-required:"true" yaml:"bucket_name" env:"GCLOUD_STORAGE_BUCKET_NAME"`
		UrlExpirationTime int    `env-required:"true" yaml:"url_expiration_time" env:"GCLOUD_STORAGE_URL_EXPIRATION_TIME"`
		UseEmulator       bool   `env-required:"true" yaml:"use_emulator" env:"GCLOUD_STORAGE_USE_EMULATOR" env-default:"false"`
		EmulatorPort      int    `env-required:"false" yaml:"emulator_port" env:"GCLOUD_STORAGE_EMULATOR_PORT"`
		UseCredentials    bool   `yaml:"use_credentials" env:"GCLOUD_STORAGE_USE_CREDENTIALS" env-default:"false"`
		AccessId          string `yaml:"access_id" env:"GCLOUD_STORAGE_ACCESS_ID"`
		PrivateKeyBase64  string `yaml:"private_key_base_64" env:"GCLOUD_STORAGE_PRIVATE_KEY_BASE_64"`
		Insecure          bool   `yaml:"insecure" env:"GCLOUD_STORAGE_INSECURE" env-default:"false"`
	}
)

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - can't read config: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

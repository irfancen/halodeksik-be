package appconfig

import (
	"halodeksik-be/app/env"
	"os"
)

var Config *AppConfig

type AppConfig struct {
	GcloudStorageCdn                        string
	GcloudStorageFolderConsultationSessions string
}

func LoadConfig() error {
	err := env.LoadEnv()
	if err != nil {
		return err
	}
	Config = &AppConfig{
		GcloudStorageCdn:                        os.Getenv("GCLOUD_STORAGE_CDN"),
		GcloudStorageFolderConsultationSessions: os.Getenv("GCLOUD_STORAGE_FOLDER_CONSULTATION_SESSIONS"),
	}
	return nil
}

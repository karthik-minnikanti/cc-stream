package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	UploadTo         string
	BucketName       string
	IPCameraURL      string
	OutputFile       string
	RecordDuration   int
	ClientSecretFile string
	TokenFilePath    string
	UploadInterval   time.Duration
}

var Env *Config

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	recordDuration, _ := strconv.Atoi(os.Getenv("RECORD_DURATION"))
	uploadInterval, _ := time.ParseDuration(os.Getenv("UPLOAD_INTERVAL"))

	Env = &Config{
		UploadTo:         os.Getenv("UPLOAD_TO"),
		BucketName:       os.Getenv("BUCKET_NAME"),
		IPCameraURL:      os.Getenv("IP_CAMERA_URL"),
		OutputFile:       os.Getenv("OUTPUT_FILE"),
		RecordDuration:   recordDuration,
		ClientSecretFile: os.Getenv("CLIENT_SECRET_FILE"),
		TokenFilePath:    os.Getenv("TOKEN_FILE_PATH"),
		UploadInterval:   uploadInterval,
	}
}

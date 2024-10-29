package main

import (
	"cc-stream/config"
	"cc-stream/handlers"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	switch config.Env.UploadTo {
	case "s3":
		for {
			fmt.Println("Recording IP camera stream...")
			if err := handlers.RecordIPCameraStream(); err != nil {
				log.Fatalf("Failed to record IP camera stream: %v", err)
			}
			fmt.Println("Recording completed.")

			if err := handlers.UploadToS3("output.mp4", *config.Env); err != nil {
				log.Fatalf("Failed to upload video to S3: %v", err)
			}

			if err := os.Remove("output.mp4"); err != nil {
				log.Printf("Failed to delete local video file: %v", err)
			} else {
				fmt.Println("Local video file deleted.")
			}

			time.Sleep(config.Env.UploadInterval)
		}

	case "youtube":
		streamURL, streamKey, err := handlers.CreateYouTubeLiveStream()
		if err != nil {
			log.Fatalf("Failed to create YouTube live stream: %v", err)
		}

		if err := handlers.StreamToYouTube(streamURL, streamKey); err != nil {
			log.Fatalf("Failed to stream to YouTube: %v", err)
		}
	}
}

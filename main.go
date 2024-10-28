package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	ipCameraURL      = "rtsp://karthik:karthik@192.168.31.224/stream1" // Replace with your IP camera URL
	outputFile       = "output.mp4"                                    // Temporary file for recording
	recordDuration   = 60                                              // Recording duration in seconds
	clientSecretFile = "client_secret.json"                            // OAuth credentials file
	tokenFilePath    = "token.json"                                    // File to store OAuth token
	uploadInterval   = 0 * time.Minute                                 // Interval to upload recordings
)

type Config struct {
	UploadTo   string `json:"upload_to"`
	BucketName string `json:"bucket_name"`
}

func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {

	config, err := loadConfig("config.json") // Path to your configuration file
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	for {
		startTime := time.Now()
		fmt.Println("Recording IP camera stream...")
		if err := recordIPCameraStream(); err != nil {
			log.Fatalf("Failed to record IP camera stream: %v", err)
		}
		fmt.Println("Recording completed.")

		// Generate video title based on recording times
		videoTitle := fmt.Sprintf("IP Camera Recording %s to %s", startTime.Format("2006-01-02 15:04:05"), startTime.Add(recordDuration*time.Second).Format("2006-01-02 15:04:05"))

		// Step 2: Upload video to YouTube

		switch config.UploadTo {
		case "s3":
			if err := uploadToS3(outputFile, *config); err != nil {
				log.Fatalf("Failed to upload video to S3: %v", err)
			}

		case "youtube":
			if err := uploadToYouTube(outputFile, videoTitle); err != nil {
				log.Fatalf("Failed to upload video to YouTube: %v", err)
			}
		}

		// Clean up the local video file after upload
		if err := os.Remove(outputFile); err != nil {
			log.Printf("Failed to delete local video file: %v", err)
		} else {
			fmt.Println("Local video file deleted.")
		}

		// Wait for the specified upload interval before the next recording
		time.Sleep(uploadInterval)
	}
}

// recordIPCameraStream records the IP camera stream using FFmpeg for the specified duration.
func recordIPCameraStream() error {
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", ipCameraURL,
		"-t", fmt.Sprintf("%d", recordDuration),
		"-c:v", "copy",
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "mp4",
		outputFile,
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start FFmpeg command: %w", err)
	}
	return cmd.Wait()
}

// uploadToYouTube uploads a video to YouTube using the YouTube Data API.
func uploadToYouTube(filename, title string) error {
	ctx := context.Background()

	// Read OAuth credentials file
	b, err := ioutil.ReadFile(clientSecretFile)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %w", err)
	}

	// Load the OAuth configuration
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	client := getClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("unable to create YouTube service: %w", err)
	}

	// Prepare the YouTube video
	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: "Recorded from an IP camera",
			Tags:        []string{"IP Camera", "Recording"},
			CategoryId:  "22", // 22 = People & Blogs
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "private"}, // Set video to private
	}

	// Open the video file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open video file: %w", err)
	}
	defer file.Close()

	// Upload the video
	call := service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("failed to upload video to YouTube: %w", err)
	}

	fmt.Printf("Video uploaded successfully! Video ID: %s\n", response.Id)
	return nil
}

// getClient retrieves an authenticated HTTP client using OAuth2 credentials.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFilePath, tok)
	}
	return config.Client(ctx, tok)
}

// tokenFromFile retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWeb requests a token from the web and returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser and enter the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to save OAuth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func uploadToS3(filename string, config Config) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"), // Specify your region
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bucket := config.BucketName // Replace with your S3 bucket name
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	fmt.Printf("File uploaded to S3: %s\n", filename)
	return nil
}

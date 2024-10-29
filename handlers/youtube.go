package handlers

import (
	"cc-stream/auth"
	"cc-stream/config"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func CreateYouTubeLiveStream() (string, string, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.Env.ClientSecretFile)
	if err != nil {
		return "", "", fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		return "", "", fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	client := auth.GetClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		return "", "", fmt.Errorf("unable to create YouTube service: %w", err)
	}

	// Create broadcast
	broadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              fmt.Sprintf("IP Camera Live Stream - %s", time.Now().Format("2006-01-02 15:04:05")),
			ScheduledStartTime: time.Now().Format(time.RFC3339),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: "private",
		},
		ContentDetails: &youtube.LiveBroadcastContentDetails{
			EnableAutoStart: true,
			EnableAutoStop:  true,
		},
	}

	broadcastInsert := service.LiveBroadcasts.Insert([]string{"snippet", "status", "contentDetails"}, broadcast)
	broadcastResponse, err := broadcastInsert.Do()
	if err != nil {
		return "", "", fmt.Errorf("error creating broadcast: %w", err)
	}

	stream := &youtube.LiveStream{
		Snippet: &youtube.LiveStreamSnippet{
			Title: "IP Camera Stream",
		},
		Cdn: &youtube.CdnSettings{
			Format:        "1080p",
			IngestionType: "rtmp",
			Resolution:    "1080p",
			FrameRate:     "30fps",
		},
	}

	streamInsert := service.LiveStreams.Insert([]string{"snippet", "cdn"}, stream)
	streamResponse, err := streamInsert.Do()
	if err != nil {
		return "", "", fmt.Errorf("error creating stream: %w", err)
	}

	bind := service.LiveBroadcasts.Bind(broadcastResponse.Id, []string{"id", "contentDetails"})
	bind.StreamId(streamResponse.Id)
	_, err = bind.Do()
	if err != nil {
		return "", "", fmt.Errorf("error binding broadcast and stream: %w", err)
	}

	return streamResponse.Cdn.IngestionInfo.IngestionAddress, streamResponse.Cdn.IngestionInfo.StreamName, nil
}

func StreamToYouTube(streamURL, streamKey string) error {
	rtmpURL := fmt.Sprintf("%s/%s", streamURL, streamKey)

	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", config.Env.IPCameraURL,
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-maxrate", "2500k",
		"-bufsize", "5000k",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ar", "44100",
		"-f", "flv",
		rtmpURL,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start FFmpeg streaming: %w", err)
	}

	return cmd.Wait()
}

func RecordIPCameraStream() error {
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", config.Env.OutputFile,
		"-t", fmt.Sprintf("%d", config.Env.UploadInterval.Abs()),
		"-c:v", "copy",
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "mp4",
		"output.mp4",
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start FFmpeg command: %w", err)
	}
	return cmd.Wait()
}

func UploadToYouTube(filename, title string) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.Env.ClientSecretFile)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
	)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	client := auth.GetClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("unable to create YouTube service: %w", err)
	}

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: "Recorded from an IP camera",
			Tags:        []string{"IP Camera", "Recording"},
			CategoryId:  "22", // 22 = People & Blogs
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "private"}, // Set video to private
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open video file: %w", err)
	}
	defer file.Close()

	call := service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("failed to upload video to YouTube: %w", err)
	}

	fmt.Printf("Video uploaded successfully! Video ID: %s\n", response.Id)
	return nil
}

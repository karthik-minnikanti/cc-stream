# IP Camera Stream Recorder and Uploader

This application records a stream from an IP camera and uploads the recording to either YouTube or Amazon S3, based on the configuration specified in `config.json`.

## Features

- **Live Streaming**: Directly streams video from an IP camera using FFmpeg without recording
- **Upload to S3**: Uploads to S3 depending on configuration
- **Flexible Configuration**: Configurable streaming settings and destinations
- **Real-time Processing**: Handles video feed in real-time without local storage

## Requirements

- [Docker](https://docs.docker.com/get-docker/)
- For YouTube uploads:
  - Google Cloud project with YouTube Data API enabled
  - OAuth 2.0 credentials (`client_secret.json`)
- For S3 uploads:
  - AWS account with S3 permissions
  - AWS credentials configured

## Installation & Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/karthik-minnikanti/cc-stream.git
   cd cc-stream
   ```

2. Configure environment variables:
   Create `.env` file in the project root:
   ```
   UPLOAD_TO=youtube          # "youtube" or "s3"
   BUCKET_NAME=youtube        # S3 bucket name (for S3 uploads)
   IP_CAMERA_URL=rtsp://user:pass@camera-ip/stream  # Your IP camera RTSP URL
   OUTPUT_FILE=output.mp4     # Local recording filename
   RECORD_DURATION=300        # Recording duration in seconds
   CLIENT_SECRET_FILE=client_secret.json  # Google OAuth credentials file
   TOKEN_FILE_PATH=token.json # YouTube auth token file
   UPLOAD_INTERVAL=5          # Interval between uploads in minutes
   ```

3. Set up credentials:
   - For YouTube: Place `client_secret.json` from Google Cloud Console in project root
   - For S3: Configure AWS credentials via AWS CLI

## Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t cc-stream .
   ```

2. Run the container:
   ```bash
   docker run -v $(pwd)/config.json:/app/config.json \
             -v $(pwd)/client_secret.json:/app/client_secret.json \
             -v $(pwd)/token.json:/app/token.json \
             cc-stream
   ```

### First Time Setup for YouTube

When running for the first time with YouTube:
1. The program will provide an authorization URL
2. Open the URL in your browser
3. Log in to your Google account and grant permissions
4. Copy the authorization code shown
5. Paste the code back into the container prompt
6. The token will be saved to token.json for future use

### Regular Usage

After initial setup, the program will:
- For S3 uploads:
  - Record from IP camera
  - Upload to S3
  - Delete local file
  - Wait for configured interval
  - Repeat process

- For YouTube uploads:
  - Create YouTube live stream
  - Stream directly from camera to YouTube

To stop the program:
- Press Ctrl+C



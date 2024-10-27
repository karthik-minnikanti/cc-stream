# IP Camera Video Recorder and YouTube Uploader

This Go application records video streams from an IP camera and uploads them to YouTube at specified intervals. The uploaded videos are named with the recording start and end times and are set to private mode.

## Features
- Records video from an IP camera using FFmpeg.
- Automatically uploads recorded videos to YouTube.
- Video files are named with the recording start and end timestamps.
- Supports private video uploads.

## Prerequisites
Before running this application, ensure you have the following installed:

- Go (version 1.16 or later)
- FFmpeg
- A Google account with access to YouTube
- OAuth 2.0 credentials for YouTube Data API

## Setup Instructions

1. **Clone the Repository**
   ```bash
   git clone git@github.com:karthik-minnikanti/cc-stream.git
   cd cc-stream
   ```

2. **Install Dependencies**
   Ensure you have the Go dependencies installed:
   ```bash
   go mod tidy
   ```

3. **Configure OAuth 2.0 Credentials**
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project.
   - Enable the YouTube Data API v3 for your project.
   - Create OAuth 2.0 credentials:
     - Go to `Credentials` > `Create Credentials` > `OAuth Client ID`.
     - Choose `Desktop App` as the application type.
     - Download the `client_secret.json` file and place it in the project root directory.

4. **Token File**
   The first time you run the application, it will prompt you to authenticate and authorize access to your YouTube account. The application will save your OAuth token in `token.json` for subsequent runs.

5. **Configure IP Camera URL**
   Update the `ipCameraURL` constant in the `main.go` file with the RTSP URL of your IP camera.

6. **Run the Application**
   Start the application using:
   ```bash
   go run main.go
   ```

## Usage
- The application will record video from the specified IP camera every 5 minutes (configurable) and upload it to YouTube.
- Each video will be named with the recording start and end times in the format:
  ```
  IP Camera Recording YYYY-MM-DD HH:MM:SS to YYYY-MM-DD HH:MM:SS
  ```

## Notes
- Make sure your IP camera supports the RTSP protocol.
- Ensure that FFmpeg is properly installed and accessible in your system's PATH.
- The YouTube Data API has quota limits; make sure to monitor your usage.
- Videos are uploaded in private mode; you can change this in the code if needed.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments
- Thanks to the developers of FFmpeg and the Go programming language for their contributions to the open-source community.

# IP Camera Stream Recorder and Uploader

This application records a stream from an IP camera and uploads the recording to either YouTube or Amazon S3, based on the configuration specified in `config.json`.

## Features

- **Record Video**: Records video from an IP camera stream using FFmpeg.
- **Upload to YouTube or S3**: Uploads recorded videos to YouTube or S3 depending on configuration.
- **Automatic Cleanup**: Deletes local video files after upload.
- **Configurable Intervals**: Set upload intervals for continuous recording and uploading.

## Requirements

- [Go](https://golang.org/doc/install)
- [FFmpeg](https://ffmpeg.org/download.html) (must be installed and accessible from the command line)
- **AWS and YouTube credentials**:
  - AWS account with S3 permissions
  - Google Cloud project with YouTube Data API enabled

## Installation

1. **Clone the Repository**:
   ```bash
   git clone git@github.com:karthik-minnikanti/cc-stream.git
   cd cc-stream
   ```

2. **Install Dependencies**:
   Use Go modules to install dependencies:
   ```bash
   go mod tidy
   ```

3. **Set up Google OAuth Credentials**:
   - Download `client_secret.json` from your Google Cloud console (with YouTube Data API enabled) and place it in the root directory.

4. **Set up AWS Credentials**:
   - Configure your AWS credentials via the AWS CLI or by setting environment variables (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`).

5. **Configuration**:
   Create a `config.json` file in the root directory to specify the upload destination and S3 bucket name:
   ```json
   {
     "upload_to": "youtube", // Set to "youtube" or "s3"
     "bucket_name": "your-s3-bucket-name" // S3 bucket name if upload_to is set to "s3"
   }
   ```

## Usage

1. **Run the Application**:
   ```bash
   go run main.go
   ```

2. **Authenticate for YouTube Uploads**:
   - The first time you run the program, it will prompt you to authenticate with your Google account. Follow the printed instructions and enter the authorization code to grant access for YouTube uploads.

3. **Set IP Camera URL**:
   - Modify the `ipCameraURL` constant in `main.go` to point to your IP camera's RTSP stream.

## Configuration Details

- **Upload Destination (`upload_to`)**:  
  Determines where the video will be uploaded. Set to `"youtube"` for YouTube or `"s3"` for Amazon S3.

- **Bucket Name (`bucket_name`)**:  
  Required if `upload_to` is set to `"s3"`. Replace `"your-s3-bucket-name"` with the name of your S3 bucket.

- **Recording Duration (`recordDuration`)**:  
  Set the duration (in seconds) for each recording session.

- **Upload Interval (`uploadInterval`)**:  
  Specifies the interval between consecutive recordings and uploads. Set to `0` for continuous recording.

## Code Structure

- **`recordIPCameraStream`**  
  Uses FFmpeg to record the IP camera stream for the specified duration and saves it as `output.mp4`.

- **`uploadToYouTube`**  
  Authenticates with the YouTube API using OAuth2 and uploads the video to YouTube.

- **`uploadToS3`**  
  Uploads the recorded video file to an Amazon S3 bucket.

## Troubleshooting

- **FFmpeg Errors**:  
  Ensure FFmpeg is installed and accessible from the command line. Use `ffmpeg -version` to check the installation.

- **Google OAuth Authentication**:  
  If authentication fails, verify that `client_secret.json` is correctly downloaded from the Google Developer Console and placed in the project root.

- **AWS S3 Upload Issues**:  
  Ensure that your AWS credentials are correctly set up and that the S3 bucket exists with appropriate permissions.

## License

This project is licensed under the MIT License.
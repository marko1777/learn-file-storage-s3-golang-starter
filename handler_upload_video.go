package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(
	w http.ResponseWriter,
	r *http.Request,
) {
	const maxMemory = 1 << 30
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory)

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't validate JWT",
			err,
		)
		return
	}

	dbVideo, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't find video",
			err,
		)
		return
	}

	if dbVideo.UserID != userID {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"video doesn't owned by user",
			err,
		)
		return
	}

	fmt.Println("uploading video", videoID, "by user", userID)

	file, header, err := r.FormFile("video")
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Unable to parse form file",
			err,
		)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))

	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	tempFile, err := os.CreateTemp("", "tubely-upload.mp4")
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Unable to create temp file",
			err,
		)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err = io.Copy(tempFile, file); err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error saving file",
			err,
		)
		return
	}

	tempFile.Seek(0, io.SeekStart)

	processedFileName, err := processVideoForFastStart(tempFile.Name())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't convert file to faststart",
			err,
		)
		return
	}

	processedFile, err := os.Open(processedFileName)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't open procesesed file",
			err,
		)
		return
	}

	defer processedFile.Close()

	aspectRatio, err := getVideoAspectRatio(processedFileName)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't get aspect ratio",
			err,
		)
		return
	}

	fileID := cfg.prefixFileID(aspectRatio, mediaType)

	_, err = cfg.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(fileID),
		Body:        processedFile,
		ContentType: aws.String(mediaType),
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Unable to put object to s3",
			err,
		)
		return
	}

	url := fmt.Sprintf("https://%s/%s", cfg.s3CfDistribution, fileID)
	fmt.Println(url)
	dbVideo.VideoURL = &url
	err = cfg.db.UpdateVideo(dbVideo)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error updateing video",
			err,
		)
		return
	}
	// dbVideo, err = cfg.dbVideoToSignedVideo(dbVideo)
	// if err != nil {
	// 	respondWithError(
	// 		w,
	// 		http.StatusInternalServerError,
	// 		"Error signing video",
	// 		err,
	// 	)
	// 	return
	// }
	respondWithJSON(w, http.StatusOK, dbVideo)
}

func (cfg *apiConfig) dbVideoToSignedVideo(
	video database.Video,
) (database.Video, error) {
	split := strings.Split(*video.VideoURL, ",")
	if len(split) != 2 {
		return database.Video{}, fmt.Errorf("url bad format")
	}

	fmt.Println(split)

	presignedURL, err := generatePresignedURL(
		cfg.s3Client,
		split[0],
		split[1],
		2*time.Minute,
	)

	if err != nil {
		return database.Video{}, err
	}

	video.VideoURL = &presignedURL

	return video, nil
}

func generatePresignedURL(
	s3Client *s3.Client,
	bucket, key string,
	expireTime time.Duration,
) (string, error) {
	presignedClient := s3.NewPresignClient(s3Client)
	req, err := presignedClient.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(expireTime),
	)

	if err != nil {
		return "", err
	}

	fmt.Println(req.URL)
	return req.URL, nil
}

func (cfg *apiConfig) prefixFileID(
	aspectRatio string,
	mediaType string,
) string {
	prefix := "other"
	switch aspectRatio {
	case "16:9":
		prefix = "landscape"
	case "9:16":
		prefix = "portrait"
	}
	fileID := getAssetPath(mediaType)
	fileID = filepath.Join(prefix, fileID)
	return fileID
}

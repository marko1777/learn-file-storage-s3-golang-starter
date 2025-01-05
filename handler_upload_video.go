package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

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
	defer os.Remove("tubely-upload.mp4")
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

	fileID := getAssetPath(mediaType)
	_, err = cfg.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(fileID),
		Body:        tempFile,
		ContentType: &mediaType,
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

	videoURL := cfg.getObjectURL(fileID)
	dbVideo.VideoURL = &videoURL
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
	respondWithJSON(w, http.StatusOK, dbVideo)
}

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		cfg.s3Bucket,
		cfg.s3Region,
		key,
	)
}

func getAssetPath(mediaType string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	ext := mediaTypeToExt(mediaType)
	id := base64.RawURLEncoding.EncodeToString(base)
	return fmt.Sprintf("%s%s", id, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	streams := struct {
		Streams []struct {
			DisplayAspectRatio string `json:"display_aspect_ratio"`
		} `json:"streams"`
	}{}
	err = json.Unmarshal(buf.Bytes(), &streams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, stream := range streams.Streams {
		if stream.DisplayAspectRatio != "" {
			return stream.DisplayAspectRatio, nil
		}
	}
	return "", fmt.Errorf("couldn't get aspect ratio")
}

func processVideoForFastStart(filePath string) (string, error) {
	processedFilePath := filePath + ".processing"
	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		processedFilePath,
	)
	buf := bytes.NewBuffer(nil)

	cmd.Stdout = buf

	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fileInfo, err := os.Stat(processedFilePath)
	if err != nil {
		return "", fmt.Errorf("could not stat processed file: %v", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("processed file is empty")
	}

	return processedFilePath, nil
}

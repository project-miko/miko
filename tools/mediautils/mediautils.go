package mediautils

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/project-miko/miko/tools/byteutils"
)

const (
	defaultTimeOut = 60 // s

	MaxFileSize = 20 * byteutils.MB
)

var (
	mimeToExt = map[string]string{
		"image/jpeg":     ".jpg",
		"image/png":      ".png",
		"image/gif":      ".gif",
		"image/svg+xml":  ".svg",
		"image/bmp":      ".bmp",
		"image/webp":     ".webp",
		"image/tiff":     ".tiff",
		"image/x-icon":   ".ico",
		"image/jp2":      ".jp2",
		"image/x-ms-bmp": ".bmp",
		"video/mp4":      ".mp4",
	}
)

type ErrMaxFileSizeExceeded struct {
	FileSize int
}

func (e *ErrMaxFileSizeExceeded) Error() string {
	return fmt.Sprintf("the download file has exceeded the maximum allowed file size, file size: %d", e.FileSize)
}

func GetExtFromMIME(mime string) (string, bool) {
	ext, ok := mimeToExt[mime]
	return ext, ok
}

func DownloadFileFromURL(url string) (string, []byte, error) {
	client := &http.Client{
		Timeout: defaultTimeOut * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	_contentLength := resp.Header.Get("Content-Length")
	contentLength, _ := strconv.Atoi(_contentLength)
	if contentLength > MaxFileSize { // do not consider the case where the Content-Length response header is not read
		return "", nil, &ErrMaxFileSizeExceeded{contentLength}
	}

	// todo read in chunks, report an error if the limit is reached
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType, data, nil
}

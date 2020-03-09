package uploader

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
)

const MaxImageSize = 10 * 1024 * 1024 // ~10MB

var allowedContentTypes = []string{"image/jpeg", "image/png"}

func ValidateURLResponse(res *http.Response, url string) (string, error) {
	contentType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if !IsImage(contentType) {
		return "", errors.New(fmt.Sprintf("Content-Type of %s is not allowed", url))
	}
	if !IsValidImageSize(res.ContentLength) {
		return "", errors.New(fmt.Sprintf("Content-Length of %s is too big", url))
	}
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", err
	}
	return ext[0], nil
}

func IsImage(contentType string) bool {
	for _, act := range allowedContentTypes {
		if contentType == act {
			return true
		}
	}
	return false
}

func IsValidImageSize(contentLength int64) bool {
	if contentLength > MaxImageSize {
		return false
	}
	return true
}

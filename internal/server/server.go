package server

import (
	"github.com/alyakimenko/image-upload/internal/uploader"
	"github.com/alyakimenko/image-upload/internal/uploader/jsonbase64uploader"
	"github.com/alyakimenko/image-upload/internal/uploader/multipartuploader"
	"github.com/alyakimenko/image-upload/internal/uploader/urluploader"
	"mime"
	"net/http"
)

type Controller struct {
	Server        *http.Server
	downloadedDir string
}

func NewController(config *Config) *Controller {
	router := http.NewServeMux()
	srv := &Controller{
		Server:        &http.Server{Addr: config.BindAddr, Handler: router},
		downloadedDir: config.DownloadedPath,
	}

	router.HandleFunc("/upload", srv.UploadImages)

	return srv
}

func (sc *Controller) UploadImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u uploader.Uploader

	switch contentType {
	case "multipart/form-data":
		u = multipartuploader.New(r, sc.downloadedDir)
	case "application/json":
		u = jsonbase64uploader.New(r, sc.downloadedDir)
	case "text/plain":
		u = urluploader.New(r, sc.downloadedDir)
	default:
		http.Error(w, "Provided Content-Type not allowed", http.StatusBadRequest)
		return
	}

	if err := u.Upload(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

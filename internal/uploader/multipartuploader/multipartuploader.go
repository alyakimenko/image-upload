package multipartuploader

import (
	"fmt"
	"github.com/alyakimenko/image-upload/internal/uploader"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

type MultipartUploader struct {
	r   *http.Request
	dir string
}

func New(r *http.Request, dir string) *MultipartUploader {
	return &MultipartUploader{r, dir}
}

func (mu *MultipartUploader) Upload() error {
	if err := mu.r.ParseMultipartForm(uploader.MaxImageSize); err != nil {
		return err
	}

	mf := mu.r.MultipartForm
	files := mf.File["files"]

	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		contentType := files[i].Header.Get("Content-Type")
		if !uploader.IsImage(contentType) {
			log.Println(fmt.Sprintf("Content-Type of %s now allowed", files[i].Filename))
			continue
		}

		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		resizedImage, err := uploader.ResizeImage(fileData)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		ext := path.Ext(files[i].Filename)
		filename := strings.TrimSuffix(files[i].Filename, ext)
		if err := uploader.SaveImage(resizedImage, mu.dir, filename, ext); err != nil {
			log.Println(err.Error())
			continue
		}

	}

	return nil
}

package jsonbase64uploader

import (
	"encoding/base64"
	"encoding/json"
	"github.com/alyakimenko/image-upload/internal/uploader"
	"io/ioutil"
	"log"
	"net/http"
)

type JSONBase64Uploader struct {
	r *http.Request
	dir string
}

func New(r *http.Request, dir string) *JSONBase64Uploader {
	return &JSONBase64Uploader{r, dir}
}

func (ju *JSONBase64Uploader) Upload() error {
	imagesData := &struct {
		Data []string `json:"data"`
	}{}

	reqBody, err := ioutil.ReadAll(ju.r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(reqBody, imagesData); err != nil {
		return err
	}

	for _, enc := range imagesData.Data {
		dec, err := base64.StdEncoding.DecodeString(enc)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		resizedImage, err := uploader.ResizeImage(dec)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		filename := uploader.GenerateFilename()
		if err := uploader.SaveImage(resizedImage, ju.dir, filename, ".jpeg"); err != nil {
			log.Println(err.Error())
			continue
		}
	}

	return nil
}

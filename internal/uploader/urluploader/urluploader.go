package urluploader

import (
	"github.com/alyakimenko/image-upload/internal/uploader"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type URLUploader struct {
	r   *http.Request
	dir string
}

func New(r *http.Request, dir string) *URLUploader {
	return &URLUploader{r, dir}
}

func (uu *URLUploader) Upload() error {
	reqBody, err := ioutil.ReadAll(uu.r.Body)
	if err != nil {
		return err
	}

	urls := strings.Split(string(reqBody), "\n")

	client := http.Client{Timeout: time.Second * 5}

	for _, url := range urls {
		res, err := client.Get(url)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer res.Body.Close()

		ext, err := uploader.ValidateURLResponse(res, url)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		resData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		resizedImage, err := uploader.ResizeImage(resData)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		filename := uploader.GenerateFilename()
		if err := uploader.SaveImage(resizedImage, uu.dir, filename, ext); err != nil {
			log.Println(err.Error())
			continue
		}

	}
	return nil
}

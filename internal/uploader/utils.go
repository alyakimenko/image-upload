package uploader

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func ResizeImage(file []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	return resize.Resize(100, 100, img, resize.Lanczos3), nil
}

func SaveImage(in image.Image, dir string, filename string, extension string) error {
	if in == nil {
		return errors.New("Cannot save " + filename)
	}

	file, err := os.Create(dir + filename + extension)
	if err != nil {
		return err
	}

	defer file.Close()

	switch extension {
	case ".jpeg", ".jpg":
		if err := jpeg.Encode(file, in, nil); err != nil {
			return err
		}
	case ".png":
		if err := png.Encode(file, in); err != nil {
			return err
		}
	}
	return nil
}

func GenerateFilename() string {
	return uuid.New().String()
}

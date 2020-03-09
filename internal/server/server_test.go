package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadImages(t *testing.T) {
	b, w := createMultipartFormData(t, "files", "./teststore/1.jpg")
	multipartContentType := w.FormDataContentType()
	log.Println(multipartContentType)

	base64JSON := struct {
		Data []string `json:"data"`
	}{
		Data: []string{
			createBase64EncodedImage(t, "./teststore/1.jpg"),
			createBase64EncodedImage(t, "./teststore/2.jpg"),
			createBase64EncodedImage(t, "./teststore/3.jpg"),
		},
	}
	base64JSONBytes, err := json.Marshal(base64JSON)
	if err != nil {
		t.Errorf("Error marshaling base64JSONBody: %v", err)
	}

	testCases := []struct {
		name         string
		contentType  string
		method       string
		expectedCode int
		body         io.Reader
	}{
		{
			name:         "urls",
			contentType:  "text/plain",
			method:       http.MethodPost,
			expectedCode: http.StatusNoContent,
			body:         bytes.NewBuffer([]byte("https://i.picsum.photos/id/237/200/300.jpg\nhttps://i.picsum.photos/id/256/200/300.jpg")),
		},
		{
			name:         "multipart/form-data",
			contentType:  multipartContentType,
			method:       http.MethodPost,
			expectedCode: http.StatusNoContent,
			body:         &b,
		},
		{
			name:         "base64 json",
			contentType:  "application/json",
			method:       http.MethodPost,
			expectedCode: http.StatusNoContent,
			body:         bytes.NewBuffer(base64JSONBytes),
		},
		{
			name:         "invalid content-type",
			contentType:  "image/jpeg",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			body:         bytes.NewBuffer([]byte("")),
		},
		{
			name:         "invalid http method",
			contentType:  "text/plain",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			body:         bytes.NewBuffer([]byte("")),
		},
	}

	config := &Config{
		BindAddr:       ":8080",
		DownloadedPath: "./downloaded/",
	}
	controller := NewController(config)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/upload", tc.body)
			req.Header.Set("Content-Type", tc.contentType)
			controller.Server.Handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var (
		b   bytes.Buffer
		fw  io.Writer
		err error
	)
	w := multipart.NewWriter(&b)

	file := mustOpen(fileName)
	defer file.Close()

	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func createBase64EncodedImage(t *testing.T, filename string) string {
	f := mustOpen(filename)
	defer f.Close()

	all, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("Error reading file content: %v", err)
	}

	return base64.StdEncoding.EncodeToString(all)
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

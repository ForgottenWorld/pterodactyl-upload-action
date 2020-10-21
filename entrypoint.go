package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var panel string
var apiKey string
var server string
var upload string

func init() {
	panel = os.Args[1]
	apiKey = os.Args[2]
	server = os.Args[3]
	upload = os.Args[4]
}

func main() {
	url, err := getSignedUrl()
	if err != nil {
		log.Fatalf("Error in signed url: %s", err)
	}

	req, err := newUploadRequest(url, upload)
	if err != nil {
		log.Fatalf("Error in upload request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req); 
	if err != nil {
		log.Fatalf("Error in main: %s", err)
	}
	defer resp.Body.Close()
}

func getSignedUrl() (string, error) {
	req, err := http.NewRequest("GET", panel+"api/client/servers/"+server+"/files/upload", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := struct {
		Attributes struct {
			Url string
		} `json:"attributes"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("%s", err)
	}
	return data.Attributes.Url, nil
}

func newUploadRequest(url string, upload string) (*http.Request, error) {
	file, err := os.Open(upload)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	form, err := writer.CreateFormFile("files", filepath.Base(upload))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(form, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

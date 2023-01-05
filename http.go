package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func EmitHttpRequest(method string, urlPath string, body []byte, writer *multipart.Writer, rsp interface{}) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(req)
	if response.StatusCode != http.StatusCreated {
		log.Printf("Request failed with response code: %d", response.StatusCode)
		return err
	}
	err = json.NewDecoder(response.Body).Decode(&rsp)
	if err != nil {
		return err
	}
	return nil
}

func EmitRequestToFileserver(filename string, input *os.File, urlPath string, rsp interface{}) error {
		// New multipart writer.
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		newFileName := fmt.Sprintf("f-%d-%s", time.Now().Unix(), filename)
		fw, err := writer.CreateFormFile("file", newFileName)
		if err != nil {
			return err
		}
		
		_, err = io.Copy(fw, input)
		if err != nil {
			return err
		}
		writer.Close()

		// 调用http请求
		err = EmitHttpRequest(http.MethodPost, urlPath, body.Bytes(), writer, rsp)

		return err
}
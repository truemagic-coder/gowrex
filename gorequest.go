package gorequest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func reqFormFile(uri string, params map[string]string, paramName, fileName string, fileBuffer *bytes.Buffer, method string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fileBuffer)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}

func reqJSON(uri string, body *bytes.Buffer, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, body)
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

// PostFormFile - POST a multipart upload with file buffer with optional params
func PostFormFile(uri string, params map[string]string, paramName, fileName string, fileBuffer *bytes.Buffer) (*http.Request, error) {
	return reqFormFile(uri, params, paramName, fileName, fileBuffer, "POST")
}

// PutFormFile - PUT a multipart upload with file buffer with optional params
func PutFormFile(uri string, params map[string]string, paramName, fileName string, fileBuffer *bytes.Buffer) (*http.Request, error) {
	return reqFormFile(uri, params, paramName, fileName, fileBuffer, "PUT")
}

// PostJSON - POST JSON to an endpoint
func PostJSON(uri string, body *bytes.Buffer) (*http.Request, error) {
	return reqJSON(uri, body, "POST")
}

// PutJSON - PUT JSON to an endpoint
func PutJSON(uri string, body *bytes.Buffer) (*http.Request, error) {
	return reqJSON(uri, body, "PUT")
}

// GetJSON - GET a JSON from an endpoint
func GetJSON(uri string) (*http.Request, error) {
	return reqJSON(uri, nil, "GET")
}

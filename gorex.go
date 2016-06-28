package gorex

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func reqFormFileDisk(r Request, params map[string]string, paramName string, filePath string) (Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return r, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		return r, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return r, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return r, err
	}
	req, err := http.NewRequest(r.Method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.Req = req
	return r, err
}

func reqFormFile(r Request, params map[string]string, paramName string, fileName string, fileBuffer *bytes.Buffer) (Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return r, err
	}
	_, err = io.Copy(part, fileBuffer)
	if err != nil {
		return r, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return r, err
	}
	req, err := http.NewRequest(r.Method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.Req = req
	return r, err
}

func reqJSON(r Request, body interface{}) (Request, error) {
	marshalled, err := json.Marshal(body)
	if err != nil {
		return r, err
	}
	jsonBuffer := bytes.NewBuffer(marshalled)
	req, err := http.NewRequest(r.Method, r.URI, jsonBuffer)
	req.Header.Add("Content-Type", "application/json")
	r.Req = req
	return r, err
}

// Request - the request object
type Request struct {
	URI     string
	Method  string
	Req     *http.Request
	Timeout time.Duration
}

// Response - the response object
type Response struct {
	Res *http.Response
	URI string
}

// FormFileDisk - a request for a multipart upload with file path with optional params
func (r Request) FormFileDisk(params map[string]string, paramName string, filePath string) (Request, error) {
	return reqFormFileDisk(r, params, paramName, filePath)
}

// FormFile - a request for a multipart upload with file buffer with optional params
func (r Request) FormFile(params map[string]string, paramName string, fileName string, fileBuffer *bytes.Buffer) (Request, error) {
	return reqFormFile(r, params, paramName, fileName, fileBuffer)
}

// JSON - a request to a JSON endpoint
func (r Request) JSON(body interface{}) (Request, error) {
	return reqJSON(r, body)
}

// Do - process the request with timeout
func (r Request) Do() (Response, error) {
	client := &http.Client{}
	client.Timeout = r.Timeout
	res, err := client.Do(r.Req)
	return Response{res, r.URI}, err
}

// JSON - decode JSON to interface
func (e Response) JSON(decoder interface{}) (*interface{}, error) {
	defer e.Res.Body.Close()
	resp := &decoder
	err := json.NewDecoder(e.Res.Body).Decode(resp)
	return resp, err
}

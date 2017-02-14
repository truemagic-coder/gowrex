package gowrex

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

// create body writer POST
func reqForm(r Request, params map[string]string, method string) (Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	defer writer.Close()
	req, err := http.NewRequest(method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if r.Headers != nil {
		for _, header := range r.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}
	if r.BasicAuth.Username != "" {
		r.Req.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}
	r.Req = req
	return r, err
}

func reqFormFileDisk(r Request, params map[string]string, paramName string, filePath string, method string) (Request, error) {
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
	req, err := http.NewRequest(method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if r.Headers != nil {
		for _, header := range r.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}
	if r.BasicAuth.Username != "" {
		r.Req.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}
	r.Req = req
	return r, err
}

func reqFormFile(r Request, params map[string]string, paramName string, fileName string, fileBuffer *bytes.Buffer, method string) (Request, error) {
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
	req, err := http.NewRequest(method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if r.Headers != nil {
		for _, header := range r.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}
	if r.BasicAuth.Username != "" {
		r.Req.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}
	r.Req = req
	return r, err
}

func reqJSON(r Request, body interface{}, method string) (Request, error) {
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, r.URI, nil)
	} else {
		marshalled, err := json.Marshal(body)
		if err != nil {
			return r, err
		}
		jsonBuffer := bytes.NewBuffer(marshalled)
		req, err = http.NewRequest(method, r.URI, jsonBuffer)
	}
	req.Header.Add("Content-Type", "application/json")
	if r.Headers != nil {
		for _, header := range r.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}
	if r.BasicAuth.Username != "" {
		r.Req.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}
	r.Req = req
	return r, err
}

func get(r Request) (Request, error) {
	req, err := http.NewRequest("GET", r.URI, nil)
	r.Req = req
	return r, err
}

// Header - a header object
type Header struct {
	Key   string
	Value string
}

// BasicAuth - a basic auth object
type BasicAuth struct {
	Username string
	Password string
}

// Request - the request object
type Request struct {
	URI       string
	Req       *http.Request
	Timeout   time.Duration
	Headers   []Header
	BasicAuth BasicAuth
}

// Response - the response object
type Response struct {
	Res *http.Response
	URI string
}

// PostForm - POST request for a multipart form data
func (r Request) PostForm(params map[string]string) (Request, error) {
	return reqForm(r, params, "POST")
}

// PutForm - PUT request for a multipart form data
func (r Request) PutForm(params map[string]string) (Request, error) {
	return reqForm(r, params, "PUT")
}

// PostFormFileDisk - POST request for a multipart upload with file path with optional params
func (r Request) PostFormFileDisk(params map[string]string, paramName string, filePath string) (Request, error) {
	return reqFormFileDisk(r, params, paramName, filePath, "POST")
}

// PutFormFileDisk - PUT request for a multipart upload with file path with optional params
func (r Request) PutFormFileDisk(params map[string]string, paramName string, filePath string) (Request, error) {
	return reqFormFileDisk(r, params, paramName, filePath, "PUT")
}

// PostFormFile - POST request for a multipart upload with file buffer with optional params
func (r Request) PostFormFile(params map[string]string, paramName string, fileName string, fileBuffer *bytes.Buffer) (Request, error) {
	return reqFormFile(r, params, paramName, fileName, fileBuffer, "POST")
}

// PutFormFile - PUT request for a multipart upload with file buffer with optional params
func (r Request) PutFormFile(params map[string]string, paramName string, fileName string, fileBuffer *bytes.Buffer) (Request, error) {
	return reqFormFile(r, params, paramName, fileName, fileBuffer, "PUT")
}

// PostJSON - POST request to a JSON endpoint
func (r Request) PostJSON(body interface{}) (Request, error) {
	return reqJSON(r, body, "POST")
}

// PutJSON - PUT request to a JSON endpoint
func (r Request) PutJSON(body interface{}) (Request, error) {
	return reqJSON(r, body, "PUT")
}

// GetJSON - GET request to a JSON endpoint
func (r Request) GetJSON() (Request, error) {
	return reqJSON(r, nil, "GET")
}

// Get - GET request to any endpoint
func (r Request) Get() (Request, error) {
	return get(r)
}

// Do - process the request with timeout
func (r Request) Do() (Response, error) {
	client := &http.Client{}
	client.Timeout = r.Timeout
	res, err := client.Do(r.Req)
	return Response{res, r.URI}, err
}

// AddHeader - add a header on the request
func (r Request) AddHeader(key string, value string) Request {
	r.Headers = append(r.Headers, Header{Key: key, Value: value})
	return r
}

// AddBasicAuth - add basic auth on the request
func (r Request) AddBasicAuth(username string, password string) Request {
	r.BasicAuth.Password = password
	r.BasicAuth.Username = username
	return r
}

// JSON - decode JSON to interface
func (e Response) JSON(decoder interface{}) (*interface{}, error) {
	defer e.Res.Body.Close()
	resp := &decoder
	err := json.NewDecoder(e.Res.Body).Decode(resp)
	return resp, err
}

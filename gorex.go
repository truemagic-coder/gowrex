package gorex

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func reqFormFile(r Request, params map[string]string, paramName, fileName string, fileBuffer *bytes.Buffer) (*http.Request, error) {
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
	req, err := http.NewRequest(r.Method, r.URI, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.Req = req
	return req, err
}

func reqJSON(r Request, body interface{}) (*http.Request, error) {
	marshalled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	jsonBuffer := bytes.NewBuffer(marshalled)
	req, err := http.NewRequest(r.Method, r.URI, jsonBuffer)
	req.Header.Add("Content-Type", "application/json")
	r.Req = req
	return req, err
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

// FormFile - a request for a multipart upload with file buffer with optional params
func (r Request) FormFile(params map[string]string, paramName, fileName string, fileBuffer *bytes.Buffer) (*http.Request, error) {
	return reqFormFile(r, params, paramName, fileName, fileBuffer)
}

// JSON - a request to a JSON endpoint
func (r Request) JSON(body interface{}, method string) (*http.Request, error) {
	return reqJSON(r, body)
}

// Do - process the request with timeout
func (r Request) Do() (*Response, error) {
	client := &http.Client{}
	client.Timeout = r.Timeout
	res, err := client.Do(r.Req)
	if err != nil {
		return nil, err
	}
	return &Response{res, r.URI}, nil
}

// JSONDecode - decode JSON to interface
func (r Response) JSONDecode(decoder interface{}) (*interface{}, error) {
	defer r.Res.Body.Close()
	resp := &decoder
	err := json.NewDecoder(r.Res.Body).Decode(resp)
	return resp, err
}

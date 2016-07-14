# Gowrex - Golang HTTP Request Library

## Why Gowrex?

Easily create multipart form/file uploads and JSON requests.

Currently this is the only public Golang library that supports multiform file uploads out of the box.

* Simple, stable, and idomatic Go API
* Multipart file & form uploading
* JSON RESTful Requests
* Supports httptest for local router testing
* Add Custom Headers
* Connection Timeout

## Install
``` bash
  go get github.com/bevanhunt/gowrex
```

## Documentation
- [GoDoc](https://godoc.org/github.com/bevanhunt/gowrex)

## Example of JSON POST with connection timeout
``` go

// JSONReceive - json response
type JSONReceive struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int64  `json:"userId"`
}

// JSONSend - json post
type JSONSend struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int64  `json:"userId"`
}

func main() {
	timeout := 10 * time.Second
	jsonData := &JSONSend{
		Title:  "fancy book",
		Body:   "this is a fancy book",
		UserID: 12,
	}
	req, err := gowrex.Request{
		URI:     "http://jsonplaceholder.typicode.com/posts",
		Timeout: timeout}.PostJSON(jsonData)
	if err != nil {
		log.Println(err)
	}
	res, err := req.Do()
	if err != nil {
		log.Println(err)
	}
	resp := &JSONReceive{}
	res.JSON(resp)
	// should print - this is a fancy book
	fmt.Println(resp.Body)
}
```
- [full source for JSON POST example](https://github.com/bevanhunt/gowrex-json-demo)


## Similar libraries
- [GoReq](https://github.com/franela/goreq)

## TODO
- tests
- more examples

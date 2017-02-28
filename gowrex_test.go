package gowrex

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("JSON", t, func() {
		// JSONReceive - json response
		type JSONReceive struct {
			ID     int64  `json:"id"`
			Title  string `json:"title"`
			Body   string `json:"body"`
			UserID int64  `json:"userId"`
		}

		// JSONSend - json post
		type JSONSend struct {
			ID     int64  `json:"id"`
			Title  string `json:"title"`
			Body   string `json:"body"`
			UserID int64  `json:"userId"`
		}

		Convey("JSON POST", func() {
			timeout := 10 * time.Second
			jsonData := &JSONSend{
				ID:     12,
				Title:  "fancy book",
				Body:   "this is a fancy book",
				UserID: 12,
			}
			req, _ := Request{
				URI:     "http://jsonplaceholder.typicode.com/posts",
				Timeout: timeout}.PostJSON(jsonData)

			res, _ := req.Do()
			resp := &JSONReceive{}
			res.JSON(resp)
			So(resp.Body, ShouldEqual, "this is a fancy book")
		})

		Convey("JSON PUT", func() {
			timeout := 10 * time.Second
			jsonData := &JSONSend{
				ID:     1,
				Title:  "fancy book",
				Body:   "this is a fancy book",
				UserID: 1,
			}
			req, _ := Request{
				URI:     "http://jsonplaceholder.typicode.com/posts/1",
				Timeout: timeout}.PutJSON(jsonData)

			res, _ := req.Do()
			resp := &JSONReceive{}
			res.JSON(resp)
			So(resp.Body, ShouldEqual, "this is a fancy book")
		})
		Convey("JSON GET", func() {
			timeout := 10 * time.Second
			req, _ := Request{
				URI:     "http://jsonplaceholder.typicode.com/posts/1",
				Timeout: timeout}.GetJSON()

			res, _ := req.Do()
			resp := &JSONReceive{}
			res.JSON(resp)
			So(resp.Body, ShouldEqual, "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto")
		})
	})
}

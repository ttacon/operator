package yolo

import (
	"github.com/robfig/pathtree"
	"net/http"
	"net/url"
	"testing"
)

func TestRouter(t *testing.T) {
	router := router{pathtree.New()}
	router.addRoute("/GET/yolo/wassup/:val/name/:name", nil)
	url, err := url.Parse("/yolo/wassup/1/name/name")
	if err != nil {
		t.Skip()
	}
	req := &http.Request{
		Method: "GET",
		URL:    url,
	}
	router.findRoute(req)
	//t.Error("having a look")
}

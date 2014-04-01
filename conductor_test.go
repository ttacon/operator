package yolo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConductor(t *testing.T) {
	c := NewConductor()
	c.Get("/yolo/:id", yoloOnlyId)

	req, err := http.NewRequest("GET", "/yolo/5", nil)
	if err != nil {
		t.Skip(err)
	}

	resp := httptest.NewRecorder()
	c.ServeHTTP(resp, req)
	resp.Flush()

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	body := resp.Body.String()
	if !strings.Contains(body, "id: 5") {
		t.Errorf("Expected body to contain \"id: 5\", got: %q", body)
	}
}

type yolo struct {
	Id int
}

func yoloOnlyId(yolo yolo, req *http.Request, resp http.ResponseWriter) {
	fmt.Println("yolo: ", yolo)
	resp.WriteHeader(http.StatusOK)
	v := fmt.Sprintf("id: %d", yolo.Id)
	fmt.Printf("val: %q\n", v)
	resp.Write([]byte(v))
}

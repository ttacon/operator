package yolo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyConductor(t *testing.T) {
	c := NewConductor()
	req, err := http.NewRequest("GET", "/yolo/5", nil)
	if err != nil {
		t.Skip(err)
	}

	resp := httptest.NewRecorder()
	c.ServeHTTP(resp, req)
	resp.Flush()

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected status of %d, got %d", http.StatusNotFound, resp.Code)
	}

	expectedBody := "404 page not found\n"
	if resp.Body.String() != expectedBody {
		t.Errorf("expected body to be %q, got %q", expectedBody, resp.Body.String())
	}
}

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
	resp.WriteHeader(http.StatusOK)
	v := fmt.Sprintf("id: %d", yolo.Id)
	resp.Write([]byte(v))
}

func TestFullBodiedPrimitiveConductor(t *testing.T) {
	c := NewConductor()
	c.Get("/yolo/:i32/:u64/:s/:f32/:b/:i8/:i16/:u8/:i/:u16/:i64/:u32/:u/:f64",
		fullBodiedPrimitiveFunction)

	req, err := http.NewRequest(
		"GET",
		"/yolo/1/10/yolo/12.0/true/1/2/8/64/15/256/3/0/-18", nil)
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
	expected := "yolo.fullBodied{I8:1, I16:2, I32:1, " +
		"I64:256, I:64, U8:0x8, U16:0xf, U32:0x3, U64:0xa, " +
		"U:0x0, S:\"yolo\", F32:12, F64:-18, B:true}"
	if body != expected {
		t.Errorf("Expected body to be %q, got: %q", expected, body)
	}
}

type fullBodied struct {
	I8  int8
	I16 int16
	I32 int32
	I64 int64
	I   int
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	U   uint
	S   string
	F32 float32
	F64 float64
	B   bool
}

func fullBodiedPrimitiveFunction(f fullBodied, req *http.Request, resp http.ResponseWriter) {
	m := fmt.Sprintf("%#v", f)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(m))
}

type struct1 struct {
	MyId int
}

type struct2 struct {
	MyId2 int
}

func twoStructsOneFunc(s1 struct1, s2 struct2, resp http.ResponseWriter) {
	resp.Write([]byte(fmt.Sprintf("%v/%v", s1, s2)))
}

func TestTwoStructs(t *testing.T) {
	c := NewConductor()
	c.Post("/twostructs/:MyId/onefunc/:MyId2", twoStructsOneFunc)

	req, err := http.NewRequest("POST", "/twostructs/2/onefunc/1", nil)
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
	expected := "{2}/{1}"
	if body != expected {
		t.Errorf("Expected body to be %q, got: %q", expected, body)
	}
}

func TestReturnAtProperTime(t *testing.T) {
	c := NewConductor()
	c.Put("/posh", a, b)

	req, err := http.NewRequest("PUT", "/posh", nil)
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
	expected := "a"
	if body != expected {
		t.Errorf("Expected body to be %q, got: %q", expected, body)
	}
}

func a(resp http.ResponseWriter) {
	resp.Write([]byte("a"))
}

func b(resp http.ResponseWriter) {
	resp.Write([]byte("b"))
}

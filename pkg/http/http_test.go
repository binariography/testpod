package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	jsonResult = `{
  "Title": "MockTest",
  "Headers": {},
  "EnvVars": [
    "TestName=Mocktest"
  ]
}`
)

func TestJSONResponse(t *testing.T) {
	srv, err := NewMockServer()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/jsontest", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	httpBody := Data{
		Title:   "MockTest",
		Headers: make(map[string]string),
		EnvVars: []string{"TestName=Mocktest"},
	}
	srv.JSONResponse(rr, req, httpBody)
	expected := jsonResult

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

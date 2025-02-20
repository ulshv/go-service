package mw

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h1"))
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h2"))
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("h3"))
	}
	chain := Chain(h1, h2, h3)
	ts := httptest.NewServer(chain)
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respBodyStr := string(respBody)
	if respBodyStr != "h1h2h3" {
		t.Error("expected data to be written in correct order - h1h2h3")
	}
}

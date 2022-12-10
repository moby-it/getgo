package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handleContainerPush(w, req)
	res := w.Result()
	_, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status ok got %v", res.StatusCode)
	}

}

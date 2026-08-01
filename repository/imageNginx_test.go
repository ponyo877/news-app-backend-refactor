package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageNginxDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/static/exists.png" {
			w.Write([]byte("png-bytes"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	repo := NewImageNginx(server.URL)

	image, err := repo.Download("exists.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(image.File) != "png-bytes" {
		t.Fatalf("file = %q, want png-bytes", image.File)
	}
	if image.Name != "exists.png" {
		t.Fatalf("name = %q", image.Name)
	}

	if _, err := repo.Download("missing.png"); err == nil {
		t.Fatal("存在しないファイルはエラーになるべき")
	}
}

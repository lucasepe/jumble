package jumble

import (
	"encoding/base64"
	"image"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadImage(t *testing.T) {
	pixel := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	data, err := base64.StdEncoding.DecodeString(pixel)
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer ts.Close()

	im, err := LoadImage(ts.URL)
	if err != nil {
		t.Error(err)
	}

	got := im.Bounds()
	want := image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
	assert.Equal(t, got, want)
}

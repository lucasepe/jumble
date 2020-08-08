package jumble

//go:generate statik -p statik -src=./assets

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	// init the embedded file system
	_ "github.com/lucasepe/jumble/statik"
	"github.com/rakyll/statik/fs"
)

// LoadImage load a image from the specified URI.
// If the URI starts with http, attempt to
// fetch the remote image with a GET verb.
// Max image size is 200 Kb.
func LoadImage(uri string) (image.Image, error) {
	const limit = 1024 * 200 // max 200 Kb

	if strings.HasPrefix(uri, "http") {
		res, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		im, _, err := image.Decode(io.LimitReader(res.Body, limit))
		return im, nil
	} else if strings.HasPrefix(uri, "assets://") {
		return LoadFromAssets(uri)
	}

	file, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(io.LimitReader(file, limit))
	return im, err
}

// LoadFromAssets load an image from the embedded filesystem.
func LoadFromAssets(uri string) (image.Image, error) {
	fn, err := imagePath(uri)
	if err != nil {
		return nil, err
	}

	sfs, err := fs.New()
	if err != nil {
		return nil, err
	}

	file, err := sfs.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	im, _, err := image.Decode(file)
	return im, err
}

func imagePath(uri string) (string, error) {
	filename := uri[len("assets://"):]
	if !strings.HasSuffix(filename, ".png") {
		filename = filename + ".png"
	}

	if strings.HasPrefix(filename, "aws_") {
		return path.Join("/aws", filename), nil
	}

	if strings.HasPrefix(filename, "azure_") {
		return path.Join("/azure", filename), nil
	}

	if strings.HasPrefix(filename, "google_") {
		return path.Join("/google", filename), nil
	}

	return "", fmt.Errorf("unknow asset: %s", filename)
}

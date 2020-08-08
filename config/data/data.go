package data

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// FetchFromURI fetch data (with limit) from an HTTP URL
func FetchFromURI(uri string, limit int64) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(io.LimitReader(response.Body, limit))
}

// FetchFromFile fetch data (with limit) from an file
func FetchFromFile(fin string, limit int64) ([]byte, error) {
	file, err := os.Open(fin)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ioutil.ReadAll(io.LimitReader(file, limit))
}

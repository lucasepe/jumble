package data

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFromURI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from scrawl!")
	}))
	defer ts.Close()

	data, err := FetchFromURI(ts.URL, 1024*10)
	if err != nil {
		t.Error(err)
	}

	want := "Hello from scrawl!"
	if got := string(data); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

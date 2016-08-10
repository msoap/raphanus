package raphanusclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HTTP(t *testing.T) {
	app := New()

	ts := getTestServer(map[string]string{
		"/test1": `response 1`,
		"/":      `404`,
	})
	defer ts.Close()

	reader, err := app.httpGet(ts.URL + "/test1")
	if err != nil {
		t.Errorf("1. httpGet() failed, error: %s", err)
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("2. httpGet() failed, error: %s", err)
	}
	if string(body) != "response 1" {
		t.Errorf("3. httpGet() failed, body: %s", string(body))
	}

	reader, err = app.httpPost(ts.URL+"/test1", nil)
	if err != nil {
		t.Errorf("1. httpPost() failed, error: %s", err)
	}
	body, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("2. httpPost() failed, error: %s", err)
	}
	if string(body) != "response 1" {
		t.Errorf("3. httpPost() failed, body: %s", string(body))
	}

	reader, err = app.httpPut(ts.URL+"/test1", nil)
	if err != nil {
		t.Errorf("1. httpPut() failed, error: %s", err)
	}
	body, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("2. httpPut() failed, error: %s", err)
	}
	if string(body) != "response 1" {
		t.Errorf("3. httpPut() failed, body: %s", string(body))
	}

	reader, err = app.httpDelete(ts.URL + "/test1")
	if err != nil {
		t.Errorf("1. httpDelete() failed, error: %s", err)
	}
	body, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("2. httpDelete() failed, error: %s", err)
	}
	if string(body) != "response 1" {
		t.Errorf("3. httpDelete() failed, body: %s", string(body))
	}

	_, err = app.httpClient("FAKE method", "http://example/test3", nil)
	if err == nil {
		t.Errorf("1. callHTTP() failed")
	}
}

func getTestServer(pathMapping map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	for path, result := range pathMapping {
		result := result
		mux.HandleFunc(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, result)
		}))
	}

	return httptest.NewServer(mux)
}

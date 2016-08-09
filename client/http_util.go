package raphanusclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// httpGet - call GET HTTP
func httpGet(URL string) (body io.ReadCloser, err error) {
	body, err = httpClient("GET", URL, nil)
	return body, err
}

// httpPost - HTTP POST call
func httpPost(URL string, bodyRequest []byte) (body io.ReadCloser, err error) {
	body, err = httpClient("POST", URL, bytes.NewReader(bodyRequest))
	return body, err
}

// httpPut - HTTP PUT call
func httpPut(URL string, bodyRequest []byte) (body io.ReadCloser, err error) {
	body, err = httpClient("PUT", URL, bytes.NewReader(bodyRequest))
	return body, err
}

// httpDelete - HTTP DELETE call
func httpDelete(URL string) (body io.ReadCloser, err error) {
	body, err = httpClient("DELETE", URL, nil)
	return body, err
}

// httpClient - call GET/POST/PUT/... by HTTP
func httpClient(HTTPMethod, URL string, bodyReq io.Reader) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	request, err := http.NewRequest(HTTPMethod, URL, bodyReq)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func httpFinalize(body io.ReadCloser) error {
	if _, err := io.Copy(ioutil.Discard, body); err != nil {
		return err
	}

	if err := body.Close(); err != nil {
		return err
	}

	return nil
}

package raphanusclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// httpGet - call GET HTTP
func (cli *Client) httpGet(URL string) (body io.ReadCloser, err error) {
	body, err = cli.httpClient("GET", URL, nil)
	return body, err
}

// httpPost - HTTP POST call
func (cli *Client) httpPost(URL string, bodyRequest io.Reader) (body io.ReadCloser, err error) {
	body, err = cli.httpClient("POST", URL, bodyRequest)
	return body, err
}

// httpPut - HTTP PUT call
func (cli *Client) httpPut(URL string, bodyRequest io.Reader) (body io.ReadCloser, err error) {
	body, err = cli.httpClient("PUT", URL, bodyRequest)
	return body, err
}

// httpDelete - HTTP DELETE call
func (cli *Client) httpDelete(URL string) (body io.ReadCloser, err error) {
	body, err = cli.httpClient("DELETE", URL, nil)
	return body, err
}

// httpClient - call GET/POST/PUT/... by HTTP
func (cli *Client) httpClient(HTTPMethod, URL string, bodyReq io.Reader) (io.ReadCloser, error) {
	request, err := http.NewRequest(HTTPMethod, URL, bodyReq)
	if err != nil {
		return nil, err
	}

	if HTTPMethod == "POST" || HTTPMethod == "PUT" {
		request.Header.Add("Content-Type", "application/json")
	}

	if len(cli.user) > 0 && len(cli.password) > 0 {
		request.SetBasicAuth(cli.user, cli.password)
	}

	response, err := cli.http.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Unauthorized")
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

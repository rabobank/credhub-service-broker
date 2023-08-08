package httpHelper

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
)

var (
	InvalidUrlError = errors.New("invalid url")
)

type HttpRequest struct {
	url        string
	header     http.Header
	tlsOptions *tls.Config
	e          error
	content    io.Reader
}

func Request(urlComponents ...string) *HttpRequest {
	result := &HttpRequest{header: http.Header{}}

	components := len(urlComponents)
	if urlComponents == nil || components == 0 {
		result.e = InvalidUrlError
		return result
	}

	url := urlComponents[0]
	lenComponent := len(url)
	if lenComponent == 0 {
		result.e = InvalidUrlError
		return result
	}
	if url[lenComponent-1] == '/' {
		url = url[:lenComponent-1]
	}

	for i := 1; i < components; i++ {
		if len(urlComponents[i]) == 0 {
			result.e = InvalidUrlError
			return result
		}
		if urlComponents[i][0] != '/' {
			url = url + "/"
		}
		url = url + urlComponents[i]
	}

	result.url = url
	return result
}

func (r *HttpRequest) WithAuthorization(authorization string) *HttpRequest {
	r.header.Set("Authorization", authorization)
	return r
}

func (r *HttpRequest) WithBearerToken(authorization string) *HttpRequest {
	r.header.Set("Authorization", "Bearer "+authorization)
	return r
}

func (r *HttpRequest) Accepting(mimeType string) *HttpRequest {
	r.header.Set("Accept", mimeType)
	return r
}

func (r *HttpRequest) Sending(mimeType string) *HttpRequest {
	r.header.Set("Content-type", mimeType)
	return r
}

func (r *HttpRequest) WithHeader(key, value string) *HttpRequest {
	r.header.Set(key, value)
	return r
}

func (r *HttpRequest) IgnoringSsl() *HttpRequest {
	r.transport().InsecureSkipVerify = true
	return r
}

func (r *HttpRequest) WithContent(content []byte) *HttpRequest {
	r.content = bytes.NewReader(content)
	return r
}

func (r *HttpRequest) DoReturningResponse(method string) (*http.Response, error) {
	client := &http.Client{}
	if r.tlsOptions != nil {
		client.Transport = &http.Transport{TLSClientConfig: r.tlsOptions}
	}

	return r.DoWithClientReturningResponse(client, method)
}

func (r *HttpRequest) DoWithClientReturningResponse(client *http.Client, method string) (*http.Response, error) {
	if r.e != nil {
		return nil, r.e
	}

	request, e := http.NewRequest(method, r.url, r.content)
	if e != nil {
		return nil, e
	}

	request.Header = r.header

	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}

	return response, nil
}

func (r *HttpRequest) Do(method string) ([]byte, error) {
	response, e := r.DoReturningResponse(method)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("failed request")
	}

	body, e := io.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	return body, nil
}

func (r *HttpRequest) DoWithClient(client *http.Client, method string) ([]byte, error) {
	response, e := r.DoWithClientReturningResponse(client, method)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("failed request")
	}

	body, e := io.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	return body, nil
}

func (r *HttpRequest) GetWithClient(client *http.Client) ([]byte, error) {
	return r.DoWithClient(client, http.MethodGet)
}

func (r *HttpRequest) Get() ([]byte, error) {
	return r.Do(http.MethodGet)
}

func (r *HttpRequest) Post() ([]byte, error) {
	return r.Do(http.MethodPost)
}

func (r *HttpRequest) PostContent(content []byte) ([]byte, error) {
	r.content = bytes.NewReader(content)
	return r.Do(http.MethodPost)
}

func (r *HttpRequest) PostContentWithClient(client *http.Client, content []byte) ([]byte, error) {
	r.content = bytes.NewReader(content)
	return r.DoWithClient(client, http.MethodPost)
}

func (r *HttpRequest) Delete() ([]byte, error) {
	return r.Do(http.MethodDelete)
}

func (r *HttpRequest) DeleteContent(content []byte) ([]byte, error) {
	r.content = bytes.NewReader(content)
	return r.Do(http.MethodDelete)
}

func (r *HttpRequest) transport() *tls.Config {
	if r.tlsOptions == nil {
		r.tlsOptions = &tls.Config{}
	}
	return r.tlsOptions
}

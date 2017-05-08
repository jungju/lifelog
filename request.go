package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type httpRequestParams struct {
	Method  string
	URL     string
	Body    interface{}
	Forms   *url.Values
	Headers map[string]string
}

type httpResponse struct {
	Body       []byte
	StatusCode int
}

func (params *httpRequestParams) request() (*httpResponse, error) {
	var req *http.Request
	var err error
	if params.Body != nil {
		jsonByte, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(params.Method, params.URL, bytes.NewBuffer(jsonByte))
	} else if params.Forms != nil {
		req, err = http.NewRequest(params.Method, params.URL, bytes.NewBufferString(params.Forms.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(params.Forms.Encode())))
	} else {
		req, err = http.NewRequest(params.Method, params.URL, nil)
	}
	if err != nil {
		return nil, err
	}
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}
	httpClient := &http.Client{Timeout: time.Duration(10 * time.Second)}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//log.Infof("HTTP request. %s %s. Response status : %d", params.Method, params.URL, resp.StatusCode)
	return &httpResponse{
		StatusCode: resp.StatusCode,
		Body:       bytes,
	}, nil
}

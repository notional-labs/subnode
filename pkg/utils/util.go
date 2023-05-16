package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func FetchUriOverHttp(url string) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := hc.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchJsonRpcOverHttp(url string, jsonBody []byte) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 10}

	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, getErr := hc.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

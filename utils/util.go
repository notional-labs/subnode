package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func FetchUriOverHttp(url string) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := hc.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchJsonRpcOverHttp(url string, jsonBody []byte) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 30}

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
		defer func() {
			_ = res.Body.Close()
		}()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func SendError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Oops! Something was wrong"))
	return err
}

func SendResult(w http.ResponseWriter, body []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	return err
}

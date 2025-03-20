package webhook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Config struct {
	url         string
	contentType string
	okCode      int
	ctx         context.Context
}

func New(url, contentType string, okCode int) *Config {
	if okCode == 0 {
		okCode = 201
	}
	if contentType == "" {
		contentType = "application/json"
	}
	return &Config{
		url:         url,
		contentType: contentType,
		okCode:      okCode,
	}
}

func (r *Config) Close() error { return nil }
func (r *Config) Open() error  { return nil }

func (r *Config) OpenWithContext(ctx context.Context) error {
	r.ctx = ctx
	return nil
}

func (r *Config) Post(data []byte) (len int, err error) {
	var ctx context.Context
	if r.ctx == nil {
		ctx = context.Background()
	} else {
		ctx = r.ctx
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", r.contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != r.okCode {
		err = errors.New(fmt.Sprintf("invalid status expected %d received %d", r.okCode, res.StatusCode))
	}

	return

}

func (r *Config) Write(data []byte) (len int, err error) {
	return r.Post(data)
}

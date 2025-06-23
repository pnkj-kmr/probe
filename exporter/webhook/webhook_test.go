package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"probe/exporter/webhook"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPost_Close(t *testing.T) {
	req := webhook.New("", "", 0)
	err := req.Close()
	assert.Nil(t, err)
}

func TestPost_Open(t *testing.T) {
	req := webhook.New("", "", 0)
	err := req.Open()
	assert.Nil(t, err)
}

func TestPost_OpenWithContext(t *testing.T) {
	req := webhook.New("", "", 0)
	err := req.OpenWithContext(context.TODO())
	assert.Nil(t, err)
}

func TestPost_Write(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := webhook.New(server.URL, "", 201)
	_, err := req.Write([]byte("abc"))
	assert.Nil(t, err)

	req = webhook.New(server.URL, "", 200)
	err = req.OpenWithContext(context.TODO())
	assert.Nil(t, err)
	_, err = req.Write([]byte("abc"))
	assert.NotNil(t, err)
}

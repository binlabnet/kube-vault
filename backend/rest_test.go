package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	srv := &Rest{
		Port:  8091,
		Ready: &atomic.Value{},
		Vault: &Vault{},
	}

	go func() {
		err := srv.Run(context.Background())
		require.NoError(t, err)
	}()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/ping", srv.Port), nil)
	require.Nil(t, err)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "must be 200")

	go func() {
		time.Sleep(100 * time.Millisecond)
		err := srv.Shutdown(context.Background())
		require.NoError(t, err)
	}()

	st := time.Now()
	assert.True(t, time.Since(st).Seconds() < 1, "should take about 100ms")

	srv2 := &Rest{
		Port:  8091,
		Ready: &atomic.Value{},
		Vault: &Vault{},
	}
	err = srv2.Run(context.Background())
	require.Error(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	srv := &Rest{
		Ready: &atomic.Value{},
		Vault: &Vault{},
	}

	go func() {
		time.Sleep(300 * time.Millisecond)
		err := srv.Shutdown(context.Background())
		require.NoError(t, err)
	}()

	err := srv.Run(context.Background())
	require.NoError(t, err)

	st := time.Now()
	assert.True(t, time.Since(st).Seconds() < 1, "should take about 100ms")
}

func TestHealthz(t *testing.T) {
	r := chi.NewRouter()

	r.Use(Healthz)

	r.Get("/t", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/healthz", nil)
	require.Nil(t, err)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "must be 200")

	req, err = http.NewRequest("GET", ts.URL+"/t", nil)
	require.Nil(t, err)
	client = &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "must be 201")
}

func TestReadyz(t *testing.T) {
	r := chi.NewRouter()

	ready := &atomic.Value{}
	ready.Store(false)
	r.Use(Readyz(ready))

	r.Get("/t", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/readyz", nil)
	require.Nil(t, err)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "must be 503")

	ready.Store(true)
	req, err = http.NewRequest("GET", ts.URL+"/t", nil)
	require.Nil(t, err)
	client = &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "must be 201")

	req, err = http.NewRequest("GET", ts.URL+"/readyz", nil)
	require.Nil(t, err)
	client = &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "must be 200")
}

func TestNotFound(t *testing.T) {
	r := chi.NewRouter()

	r.NotFound(NotFound)

	r.Get("/t", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/notFound", nil)
	require.Nil(t, err)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "must be 404")

	req, err = http.NewRequest("GET", ts.URL+"/t", nil)
	require.Nil(t, err)
	client = &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "must be 201")
}

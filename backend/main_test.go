package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestApp_Run(t *testing.T) {
	os.Args = []string{"", "--debug", "--port=8090"}
	done := make(chan bool)
	go func() {
		time.Sleep(800 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.Nil(t, err)
	}()
	go func() {
		st := time.Now()
		main()
		assert.True(t, time.Since(st).Seconds() < 1, "should take about 500msec")
		done <- true
	}()
	time.Sleep(500 * time.Millisecond)
	resp, err := http.Get("http://127.0.0.1:8090/ping")
	require.Nil(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	assert.Equal(t, 200, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, ".", string(body))
	<-done
}

package main

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	go main()
	// Wait for server to start
	for {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", "8080"), time.Millisecond*1000)
		if conn != nil {
			conn.Close()
			break
		}
	}
	res, err := http.Get("http://localhost:8080")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, resBody)
	assert.NotContains(t, string(resBody), "<title>An error occured!</title>")
}

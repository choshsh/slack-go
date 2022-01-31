package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var url = "slack weebhook url"

func TestNewSlack(t *testing.T) {
	slack := NewSlackSender(url, false)
	assert.NotEqual(t, slack, nil)
}

func TestSetPayload(t *testing.T) {
	slack := NewSlackSender(url, false)
	slack.SetPayload("test")
	assert.NotEqual(t, len(slack.Payload.Blocks), 0)
}

func TestUrl(t *testing.T) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, resp.StatusCode, 400)
}

func TestSend(t *testing.T) {
	slack := NewSlackSender(url, false)
	slack.SetPayload("test")
	statusCode := slack.Send()
	assert.Equal(t, statusCode, http.StatusOK)
}

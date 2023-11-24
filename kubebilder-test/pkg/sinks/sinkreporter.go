package sinks

import (
	"net/http"
	"test.kubebuilder.io/project/api/v1alpha1"
	"time"
)

type ISink interface {
	Configure(config v1alpha1.Kubegpt, c Client)
	Emit(results v1alpha1.ResultSpec) error
}

func NewSink(sinkType string) ISink {
	switch sinkType {
	case "slack":
		return &SlackSink{}
		// 추가적인 Sink Providers
	default:
		// 오류 반환
		return &SlackSink{}
	}
}

type Client struct {
	hclient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		hclient: client,
	}
}

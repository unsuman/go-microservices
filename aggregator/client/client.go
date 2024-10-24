package client

import (
	"github.com/unsuman/go-microservices/types"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateDistance(d types.Distance) error {
	// httpc := http.DefaultClient
	return nil
}

package Helper

import (
	"github.com/nats-io/go-nats-streaming"
)

type NatsClient struct {
	ClusterID   string
	ClientID    string
	ShowTime    bool
	StartSeq    uint64
	StartDelta  string
	DeliverAll  bool
	DeliverLast bool
	Durable     string
	Qgroup      string
	Unsubscribe bool
	URL         string
	Topic       string
}

var Config = []NatsClient{
	NatsClient{
		URL:       stan.DefaultNatsURL,
		ClusterID: "test-cluster",
		ShowTime:  false,
		// Subscription options
		StartSeq:    0,
		DeliverAll:  false,
		DeliverLast: true,
		StartDelta:  "",
		Durable:     "",
		Qgroup:      "testgroup",
		Unsubscribe: false,
		Topic:       "get_data",
	},
}

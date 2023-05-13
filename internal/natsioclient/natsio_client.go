package natsioclient

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsIOClient struct {
	conn        *nats.Conn
	topic       string
	EncodedConn *nats.EncodedConn
}

func (c *NatsIOClient) LogConnInfo() {
	log.Printf("Connected to: Name: %v, ID: %v\n", c.conn.ConnectedClusterName(), c.conn.ConnectedServerName())
	log.Printf("Server: %v\n", c.conn.ConnectedAddr())
	log.Printf("Cluster: %v\n", c.conn.ConnectedUrl())
	for _, s := range c.conn.Servers() {
		log.Printf("Servers: %v\n", s)
	}
}

func (c *NatsIOClient) Close() {
	c.conn.Flush()
	c.conn.Close()
}

func (c *NatsIOClient) IsConnected() bool {
	return c.conn.IsConnected()
}

// getEncodedConn returns the underlying nats.EncodedConn
func (c *NatsIOClient) GetEncodedConn() *nats.EncodedConn {
	return c.EncodedConn
}

func NewNatsIOClient(endpoint string, topic string) (*NatsIOClient, error) {
	opts := nats.Options{
		RetryOnFailedConnect: true,
		AllowReconnect:       true,
		MaxReconnect:         -1,
		ReconnectWait:        5 * time.Second,
		Timeout:              1 * time.Second,
		Servers:              []string{endpoint},
	}
	nc, err := opts.Connect()

	NatsIOClient := &NatsIOClient{
		conn:  nc,
		topic: topic,
	}

	nc.SetReconnectHandler(func(nc *nats.Conn) {
		// Note that this will be invoked for the first asynchronous connect.
		NatsIOClient.LogConnInfo()
	})

	if err != nil {
		log.Printf("error creating initial connection: %v\n", err)
		return NatsIOClient, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Printf("error upgrading to encoded connection: %v\n", err)
		return NatsIOClient, err
	}

	NatsIOClient.EncodedConn = ec

	return NatsIOClient, nil
}

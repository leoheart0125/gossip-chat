package gossip

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
)

const GossipTopic = "gossip-chat"

type PubSub struct {
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

func SetupGossipChat(ctx context.Context, h host.Host) (*PubSub, error) {
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create gossipsub: %w", err)
	}

	topic, err := ps.Join(GossipTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to join topic: %w", err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		topic.Close()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		sub.Cancel()
		topic.Close()
	}()

	return &PubSub{
		ps:    ps,
		topic: topic,
		sub:   sub,
	}, nil
}

// Publish sends a message to the gossip topic
func (gc *PubSub) Publish(ctx context.Context, msg []byte) error {
	return gc.topic.Publish(ctx, msg)
}

// Messages returns a channel of incoming messages
func (gc *PubSub) Messages(ctx context.Context) (*pubsub.Message, error) {
	// Use a select to catch both messages and context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return gc.sub.Next(ctx)
	}
}

// Close cleanups resources
func (gc *PubSub) Close() error {
	gc.sub.Cancel()
	return gc.topic.Close()
}

package msg

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/pubsub-util/pkg/models"
)

type PubSubHandler struct {
	ctx    context.Context
	client *pubsub.Client
}

func NewPubSubHandler(project string) (*PubSubHandler, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Err(err, "could not init pubsub client")
		return nil, err
	}
	return &PubSubHandler{client: client, ctx: ctx}, nil
}

func (h *PubSubHandler) CleanUp() {
	h.client.Close()
}

// func NewPubSubHandlerMock() (*PubSubHandler, error) {
// 	ctx := context.Background()
// 	client, err := pubsub.NewClient()
// 	if err != nil {
// 		log.Err(err, "could not init pubsub client")
// 		return nil, err
// 	}
// 	return &PubSubHandler{client: client, ctx: ctx}, nil
// }

func (h *PubSubHandler) PublishSubscribe(msg models.PubSubMessage, topicName, subscriptionName string, timeout time.Duration) error {
	err := h.Publish(topicName, msg)
	if err != nil {
		return err
	}

	data, _, err := h.Subscribe(subscriptionName, timeout)
	if err != nil {
		return err
	}
	log.Info(data)
	return nil
}

func (h *PubSubHandler) Publish(topicName string, msg models.PubSubMessage) error {
	// Publish "hello world" on topic1.
	topic := h.client.Topic(topicName)
	res := topic.Publish(h.ctx, &pubsub.Message{
		Data:       msg.Data,
		Attributes: msg.Attr,
	})
	// The publish happens asynchronously.
	// Later, you can get the result from res:
	msgId, err := res.Get(h.ctx)
	if err != nil {
		log.Err(err, "Could not get pub sub message result for topic", topicName)
		return err
	}
	log.Info("Sent message with id", msgId, "to", h.client.Project()+"/"+topicName)
	return nil
}

func (h *PubSubHandler) Subscribe(subscriptionName string, timeout time.Duration) ([]byte, map[string]string, error) {
	log.Info("Subscribe to", subscriptionName)
	var message *pubsub.Message

	// Use a callback to receive messages via subscription1.
	sub := h.client.Subscription(subscriptionName)
	sub.ReceiveSettings.MaxOutstandingMessages = 5
	sub.ReceiveSettings.Synchronous = true
	cCxt, _ := context.WithTimeout(h.ctx, time.Second*timeout)
	err := sub.Receive(cCxt, func(ctx context.Context, m *pubsub.Message) {
		message = m
		m.Ack() // Acknowledge that we've consumed the message.
		// cancel()
	})
	if err != nil {
		log.Err(err, "Receive pubsub message error")
		return nil, nil, err
	}

	if message == nil {
		log.Info("Received pubsub message empty, might be a timeout.")
		err := errors.New("empty pub sub message")
		return nil, nil, err
	}
	return message.Data, message.Attributes, nil
}

func (h *PubSubHandler) SubscribeHandle(subscriptionName string, handler func(ctx context.Context, m *pubsub.Message)) {
	log.Info("Subscribe to", subscriptionName)

	// Use a callback to receive messages via subscription1.
	sub := h.client.Subscription(subscriptionName)
	sub.ReceiveSettings.MaxOutstandingMessages = 5
	sub.ReceiveSettings.Synchronous = true
	cCxt, _ := context.WithCancel(h.ctx)
	err := sub.Receive(cCxt, handler)
	if err != nil {
		log.Err(err, "Receive pubsub message error")
		return
	}
}

package models

// pub sub

type PubSubMessage struct {
	Data []byte
	Attr map[string]string
}

type PubSubEnvelope struct {
	Message      PubSubMessage `json:"message,omitempty"`
	Subscription string        `json:"subscription,omitempty"`
}

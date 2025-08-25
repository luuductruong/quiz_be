package pubsub

// Topic represents a messaging destination where messages can be published.
// It provides a String method to get the topic's string representation.
type Topic interface {
	String() string
}

// Publisher defines an interface for publishing messages to topics.
// It provides methods for publishing both Message objects and raw byte data.
type Publisher interface {
	// Publish sends a Message to the specified Topic.
	Publish(topic Topic, msg *Message) error
	// PublishRaw sends raw byte data to the specified Topic.
	PublishRaw(topic Topic, msg []byte) error
	// Topic returns the string representation of a topic given its ID.
	Topic(id string) Topic
}

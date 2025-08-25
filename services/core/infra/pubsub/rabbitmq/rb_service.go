/*
Package rabbitmq provides RabbitMQ implementation of the pubsub interfaces

	for message publishing and subscription handling.
*/
package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// RbConfig holds the configuration parameters for RabbitMQ connection and message routing.
type RbConfig struct {
	// Host is the RabbitMQ server hostname
	Host string `config:"host"`
	// Port is the RabbitMQ server port
	Port string `config:"port"`
	// Username for RabbitMQ authentication
	Username string `config:"username"`
	// Password for RabbitMQ authentication
	Password string `config:"password"`
	// ExchangeName is the name of the exchange to declare
	ExchangeName string `config:"exchange_name"`
	// Topic is the name of the queue to declare
	Topic string `config:"topic"`
	// Subscribes is a list of routing keys to bind to the queue
	Subscribes []string `config:"subscribes"`
}

// Service represents a RabbitMQ connection handler that manages the connection,
// channel, and message routing configuration.
type Service struct {
	URL            string           // RabbitMQ connection string
	Connection     *amqp.Connection // The active RabbitMQ connection
	Channel        *amqp.Channel    // The active RabbitMQ channel
	QueueName      string           // Name of the declared queue
	ExchangeName   string           // Name of the declared exchange
	ExchangeType   string           // Exchange routing behavior
	ReconnectDelay time.Duration    // Time to wait between reconnection attempts
	Subscribes     []string         // List of routing keys for queue binding
}

/*
NewService creates a new RabbitMQ service instance with the provided configuration.

	It sets up the connection URL and default parameters for the service.
*/
func NewService(config *RbConfig) *Service {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
	fmt.Println("exchange name: ", config.ExchangeName, " topic: ", config.Topic, " subscribes: ", config.Subscribes, " url: ", url, "")
	return &Service{
		URL:            url,
		ExchangeName:   config.ExchangeName,
		ExchangeType:   "direct",
		QueueName:      config.Topic,
		ReconnectDelay: time.Second * 5,
		Subscribes:     config.Subscribes,
	}
}

/*
Connect establishes a connection to the RabbitMQ server and sets up the channel,
exchange, queue, and bindings according to the service configuration.

It performs the following steps:

1. Establishes connection to RabbitMQ server

2. Creates a channel for communication

3. Declares exchange if ExchangeType is specified

4. Declares queue for message consumption

5. Binds queue to exchange with routing keys

Returns:
  - nil if connection and setup are successful
  - error describing the failure if any operation fails during setup
*/
func (c *Service) Connect() error {
	var err error
	/* Step 1: Connect to RabbitMQ server using the configured URL */
	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return err
	}

	/* Step 2: Create a channel for communication */
	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return err
	}

	/* Step 3: Declare exchange if type is specified */
	if c.ExchangeType != "" {
		err = c.Channel.ExchangeDeclare(
			c.ExchangeName, // name of the exchange
			c.ExchangeType, // type of exchange (direct)
			true,           // durable: survive broker restart
			false,          // auto-deleted when not in use
			false,          // internal: not used by clients
			false,          // no-wait: wait for confirmation
			nil,            // arguments
		)
		if err != nil {
			return err
		}
	}

	// Step 4: Declare queue for consuming messages
	_, err = c.Channel.QueueDeclare(
		c.QueueName, // name of the queue
		true,        // durable: survive broker restart
		false,       // delete when unused
		false,       // exclusive: used by only one connection
		false,       // no-wait: wait for confirmation
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	// Step 5: Bind queue to exchange with routing keys
	for _, subscribe := range c.Subscribes {
		err = c.Channel.QueueBind(
			c.QueueName,    // queue name
			subscribe,      // routing key
			c.ExchangeName, // exchange name
			false,          // no-wait: wait for confirmation
			nil,            // arguments
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close will close the connection and channel to the RabbitMQ server.
func (c *Service) Close() error {
	var err error
	if c.Channel != nil {
		err = c.Channel.Close()
		if err != nil {
			return err
		}
	}

	if c.Connection != nil {
		err = c.Connection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
Publish sends a message to the specified topic through the RabbitMQ exchange.
If the publish operation fails due to connection issues, it attempts to reconnect
and retry the publish operation at once.

Parameters:
  - message: The message content as a byte slice
  - topic: The routing key for message delivery

Returns:
  - nil if a message is published successfully on the first attempt or after reconnect
  - error if publishing fails after a single reconnection and a retry attempt
*/
func (c *Service) Publish(message []byte, topic string) error {
	// Step 1: Attempt to publish the message to RabbitMQ exchange
	err := c.Channel.PublishWithContext(context.Background(), // Use default context
		c.ExchangeName, // The exchange to publish to
		topic,          // Routing key for message delivery
		false,          // mandatory: don't return message if no queue bound
		false,          // immediate: don't return message if no consumer available
		amqp.Publishing{
			ContentType: "text/plain", // Set message content type
			Body:        message,      // Actual message content
		},
	)

	// Step 2: Handle publish failure by attempting reconnection
	if err != nil {
		// Try to re-establish connection if publish failed
		err = c.Reconnect()

		// Step 3: Retry publishing after reconnection
		err = c.Channel.PublishWithContext(context.Background(),
			c.ExchangeName,
			topic,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        message,
			},
		)
	}
	return err
}

/*
Reconnect attempts to re-establish the connection to the RabbitMQ server.
It will continuously retry the connection using the configured ReconnectDelay
until a successful connection is established.

Returns:
  - nil when connection is successfully established
  - This method will not return an error as it continuously retries until successful
*/
func (c *Service) Reconnect() error {
	for {
		err := c.Connect()
		if err == nil {
			return nil
		}
		log.Printf("Failed to connect to RabbitMQ server. Retrying in %d seconds\n", c.ReconnectDelay)
		time.Sleep(c.ReconnectDelay * time.Second)
	}
}

/*
Consume sets up message consumption from the configured queue.
It creates a channel for receiving messages that can be processed by the application.

Returns:
  - A channel of amqp.Delivery for receiving messages
  - Error if setting up consumption fails
*/
func (c *Service) Consume() (<-chan amqp.Delivery, error) {
	// Set up consumption from the queue with the following parameters:
	// queue - name of the queue to consume from
	// consumer - consumer identifier (empty string for auto-generation)
	// autoAck - false to require explicit message acknowledgment
	// exclusive - false to allow multiple consumers on the queue
	// noLocal - false to receive messages published by this connection
	// noWait - false to wait for server confirmation
	// args - no additional arguments needed
	return c.Channel.Consume(
		c.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

/*
ListenAndServe starts a continuous message consumption loop that processes
incoming messages using the provided callback function. It handles connection
failures by attempting to reconnect automatically.

Parameters:
  - callback: Function that processes received messages

The function will:
 1. Set up message consumption
 2. Process incoming messages using callback
 3. Automatically attempt to reconnect if the connection is lost
 4. Acknowledge messages after successful processing (Ack)
 5. Continue processing on callback errors without negative acknowledgment (no Nack)

Returns error if initial consumption setup fails or an unrecoverable error occurs.
*/
func (c *Service) ListenAndServe(callback func([]byte) error) error {
	// Initialize the message delivery channel for consuming messages from the queue
	// This channel will receive all messages published to the bound queue
	deliveryChan, err := c.Consume()
	if err != nil {
		return err
	}

	// Track connection status to handle reconnection logic
	connected := true

	// Infinite loop to continuously process messages
	for {
		// Wait for messages from the delivery channel
		select {
		case delivery, ok := <-deliveryChan:
			// Check if the channel is closed (connection lost)
			if !ok {
				connected = false
				log.Println("Connection closed. Attempting to reconnect...")
				// Keep trying to reconnect until successful
				for !connected {
					err := c.Reconnect()
					if err != nil {
						log.Printf("Reconnect failed: %v. Retrying in %d seconds...\n", err, c.ReconnectDelay)
						time.Sleep(c.ReconnectDelay * time.Second)
					} else {
						log.Println("Reconnected successfully.")
						connected = true
					}
				}
				// Re-establish the delivery channel after reconnection
				deliveryChan, err = c.Consume()
				if err != nil {
					return err
				}
				continue
			}

			// Process the received message using the provided callback function
			err := callback(delivery.Body)
			if err != nil {
				//delivery.Nack(false, true)
				log.Printf("Message processing error: %v\n", err)
				//continue
			}
			// Acknowledge the message after processing, regardless of callback errors
			// false parameter means we're only acknowledging this single delivery
			delivery.Ack(false)
		}
	}
}

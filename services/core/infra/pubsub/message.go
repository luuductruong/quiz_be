package pubsub

//
//import (
//	"cloud.google.com/go/pubsub"
//	"context"
//	"encoding/json"
//	"errors"
//	"github.com/aws/aws-sdk-go/service/sns"
//	"github.com/aws/aws-sdk-go/service/sqs"
//	"github.com/quiz_be/services/core/infra/logger"
//	"google.golang.org/api/option"
//	"time"
//)
//
//const (
//	// googleCredentialDefault - Using default GCP credential.
//	googleCredentialDefault = "Default"
//	// googleCredentialFile - Using credential specified in file.
//	googleCredentialFile = "File"
//	awsPubsubKind        = "aws_pubsub"
//	googlePubsubKind     = "google_pubsub"
//	redisPubsubKind      = "redis_pubsub"
//	RbmqPubsubKind       = "rbmq_pubsub"
//)
//
//type UUID = string
//
//// Message struct.
//// This is whole message payload that will be send to the bus.
//type Message struct {
//	// Message identifier. The original ID will be set in case of retry for a (dead) message.
//	ID         UUID `json:"id"`
//	OriginalID UUID `json:"original_id"`
//
//	// Context of the message including entity ID, ecosystem ID, trace ID, access token and checksum.
//	EntityID UUID   `json:"entity_id"`
//	TraceID  string `json:"trace_id"`
//	Token    string `json:"token"`
//
//	// Name of the message, corresponding payload and its checksum
//	Name    MessageName `json:"name"`
//	Payload []byte      `json:"payload"`
//
//	// Metadata of the message. Retry count is the number of retry to handle the message.
//	RetryCount int       `json:"retry_count"`
//	CreatedAt  time.Time `json:"created_at"`
//
//	Error error `json:"error"`
//}
//
//func (m Message) ScanPayload(payload interface{}) error {
//	return json.Unmarshal(m.Payload, payload)
//}
//
//// SubscriptionHandler func
//type SubscriptionHandler func(context.Context, *Message) error
//
//type SubscriptionRouter = func(msg *Message) SubscriptionHandler
//
//// SubscriptionInterceptor func
//type SubscriptionInterceptor func(ctx context.Context, info *Message, next SubscriptionHandler) error
//
//func RoutingHandler(route SubscriptionRouter) SubscriptionHandler {
//	return func(ctx context.Context, msg *Message) error {
//		fn := route(msg)
//		return fn(ctx, msg)
//	}
//}
//
//// PubSub interface
//type PubSub interface {
//	Publisher() Publisher
//	Subscriber() Subscriber
//	Close()
//}
//
//// Topic interface
//type Topic interface {
//	String() string
//}
//
//// Publisher do publish the message to queue
//type Publisher interface {
//	Publish(topic Topic, msg *Message) error
//	PublishRaw(topic Topic, msg []byte) error
//	Topic(id string) Topic
//}
//
//// Subscription string
//type Subscription interface {
//	String() string
//}
//
//// Subscriber do pulling message from queue
//type Subscriber interface {
//	Subscription(id string) Subscription
//	Unsubscribe(Subscription)
//	Subscribe(subscription Subscription, handler SubscriptionHandler)
//	Use(middlewares ...SubscriptionInterceptor)
//	StartReceiving() error
//	StopReceiving()
//}
//
//// Config config
//type Config struct {
//	Kind         string        `config:"kind"`
//	GooglePubSub *GoogleConfig `config:"google_pubsub"`
//	//RedisPubSub  *redis.Config      `config:"redis_pubsub"`
//	//RbmqPubsub   *RbConfig          `config:"rbmq_pubsub"`
//	//AwsPubsub    *awssession.Config `config:"aws_pubsub"`
//	Subscription string   `config:"subscription"`
//	Topic        string   `config:"topic"`
//	Subscribes   []string `config:"subscribes"`
//}
//
//// GoogleConfig struct
//type GoogleConfig struct {
//	ProjectID      string `config:"project_id"`
//	CredentialType string `config:"credential_type"`
//	CredentialPath string `config:"credential_path"`
//}
//
//// AwsConfig struct
//type AwsConfig struct {
//	ProjectID       string `config:"project_id"`
//	Region          string `config:"region"`
//	AccessKeyID     string `config:"access_key_id"`
//	SecretAccessKey string `config:"secret_access_key"`
//	AccessToken     string `config:"access_token"`
//	QueueName       string `config:"queue_name"`
//}
//
//// NewPubSub instance
//func NewPubSub(logger logger.Logger, psConfig *Config) (PubSub, error) {
//	switch psConfig.Kind {
//	case googlePubsubKind:
//		var ggClient *pubsub.Client
//		var err error
//		if ggps := psConfig.GooglePubSub; ggps.CredentialType == googleCredentialDefault {
//			ggClient, err = pubsub.NewClient(context.Background(), ggps.ProjectID)
//			if err != nil {
//				return nil, err
//			}
//		} else if ggps.CredentialType == googleCredentialFile {
//			options := option.WithCredentialsFile(psConfig.GooglePubSub.CredentialPath)
//			ggClient, err = pubsub.NewClient(context.Background(), psConfig.GooglePubSub.ProjectID, options)
//			if err != nil {
//				return nil, err
//			}
//		} else {
//			return nil, errors.New("not supported credential type")
//		}
//		return &ggPs{
//			client:     ggClient,
//			publisher:  newGgPublisher(logger, ggClient),
//			subscriber: newGgSubscriber(logger, ggClient),
//		}, nil
//	case awsPubsubKind:
//		awsSess, err := awssession.New(psConfig.AwsPubsub)
//		if err != nil {
//			return nil, err
//		}
//
//		return &awsPs{
//			publisher:  newAwsPublisher(logger, sns.New(awsSess)),
//			subscriber: newAwsSubscriber(logger, sqs.New(awsSess)),
//		}, nil
//	case RbmqPubsubKind:
//		psConfig.RbmqPubsub.Topic = psConfig.Topic
//		psConfig.RbmqPubsub.Subscribes = psConfig.Subscribes
//		return NewRbPubSub(logger, psConfig.RbmqPubsub), nil
//	default:
//		return nil, errors.New("not supported pubsub kind")
//	}
//}

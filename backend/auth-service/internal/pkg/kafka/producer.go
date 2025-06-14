package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/russian-steam/auth-service/internal/config"
)

// EventPublisher defines the interface for publishing events.
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, data interface{}) error // Changed topic to eventType
	Close()
}

// CloudEvent defines the structure for our domain events, adhering to CloudEvents spec.
// As per project_api_standards.md (section 3.1)
type CloudEvent struct {
	SpecVersion     string      `json:"specversion"`
	ID              string      `json:"id"`
	Source          string      `json:"source"` // e.g., "platform.auth-service"
	Type            string      `json:"type"`   // e.g., "com.platform.auth.user.registered.v1"
	Time            time.Time   `json:"time"`
	DataContentType string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
	TraceID         string      `json:"traceid,omitempty"`
	CorrelationID   string      `json:"correlationid,omitempty"`
}

// KafkaProducer implements EventPublisher for Kafka.
type KafkaProducer struct {
	producer    *ckafka.Producer
	topicPrefix string
	serviceName string // e.g. "platform.auth-service"
}

// NewKafkaProducer creates a new Kafka producer.
// If Kafka brokers are not configured, it returns a NoOpProducer.
func NewKafkaProducer(cfg *config.KafkaConfig, serviceName string) (EventPublisher, error) {
	if cfg.Brokers == "" {
		log.Println("Kafka brokers not configured, Kafka producer disabled. Using NoOpProducer.")
		return &NoOpProducer{}, nil
	}
	p, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"acks":              "all", // Wait for all in-sync replicas to ack
		"retries":           3,
		"retry.backoff.ms":  1000,
		// "delivery.timeout.ms": 5000, // Renamed from message.timeout.ms for librdkafka >= v1.6.0
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Println("Kafka producer created successfully.")

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *ckafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Kafka delivery failed for message to topic %s: %v\n", *ev.TopicPartition.Topic, ev.TopicPartition.Error)
				} else {
					log.Printf("Kafka message delivered to %s [%d] at offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			case ckafka.Error:
				log.Printf("Kafka producer error: %v\n", ev)
			}
		}
	}()

	return &KafkaProducer{producer: p, topicPrefix: cfg.TopicPrefix, serviceName: serviceName}, nil
}

func (kp *KafkaProducer) Publish(ctx context.Context, eventType string, data interface{}) error {
	if kp.producer == nil {
		return fmt.Errorf("kafka producer is not initialized")
	}
	// Construct topic from eventType, e.g., "com.platform.auth.user.registered.v1" -> "com.platform.auth.user.events.v1"
	// This assumes a convention where specific event types map to a more general event topic for that entity.
	parts := strings.Split(eventType, ".")
	if len(parts) < 4 { // Expecting at least service.entity.verb.version
		return fmt.Errorf("invalid eventType format: %s, expected e.g. com.platform.auth.user.registered.v1", eventType)
	}
	// Example: com.platform.auth.user.registered.v1 -> topic com.platform.auth.user.events.v1
	// Example: com.platform.auth.email.verification_requested.v1 -> topic com.platform.auth.email.events.v1
	// Example: com.platform.auth.session.revoked.v1 -> topic com.platform.auth.session.events.v1
	entityTopicName := strings.Join(parts[0:len(parts)-2], ".") // e.g., com.platform.auth.user
	topic := kp.topicPrefix + entityTopicName + ".events.v1"

	event := CloudEvent{
		SpecVersion:     "1.0",
		ID:              uuid.NewString(),
		Source:          kp.serviceName,
		Type:            eventType,
		Time:            time.Now().UTC(),
		DataContentType: "application/json",
		Data:            data,
	}

	// Extract traceID and correlationID from context if available (example)
	// if traceID, ok := ctx.Value("traceIDKey").(string); ok { event.TraceID = traceID }
	// if correlationID, ok := ctx.Value("correlationIDKey").(string); ok { event.CorrelationID = correlationID }

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal CloudEvent: %w", err)
	}

	deliveryChan := make(chan ckafka.Event)
	// No defer close(deliveryChan) here as Produce sends on it and then closes it internally (or the event loop does)

	err = kp.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          payload,
		Key:            []byte(event.ID),
	}, deliveryChan)

	if err != nil {
		// If an error occurs during Produce call itself (e.g., queue full)
		close(deliveryChan) // Ensure channel is closed if Produce fails before sending
		return fmt.Errorf("kafka produce error: %w", err)
	}

	// Wait for delivery report with timeout
	select {
	case e := <-deliveryChan:
		m := e.(*ckafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("kafka delivery failed: %w", m.TopicPartition.Error)
		}
		// Successfully produced
	case <-ctx.Done():
		return ctx.Err() // Context cancelled or timed out
	case <-time.After(10 * time.Second): // Timeout for Produce operation including delivery report
		return fmt.Errorf("kafka produce timed out for event type %s to topic %s", eventType, topic)
	}
	return nil
}

func (kp *KafkaProducer) Close() {
	if kp.producer != nil {
		log.Println("Flushing and closing Kafka producer...")
		remaining := kp.producer.Flush(15 * 1000) // 15 seconds timeout
		if remaining > 0 {
			log.Printf("WARN: %d messages still in Kafka producer queue after flush timeout", remaining)
		}
		kp.producer.Close()
		log.Println("Kafka producer closed.")
	}
}

// NoOpProducer is an EventPublisher that does nothing but logs.
type NoOpProducer struct{}

func (p *NoOpProducer) Publish(ctx context.Context, eventType string, data interface{}) error {
	// Reconstruct topic for logging, similar to how KafkaProducer does it
	parts := strings.Split(eventType, ".")
	topic := "unknown-topic"
	if len(parts) >= 4 {
		entityTopicName := strings.Join(parts[0:len(parts)-2], ".")
		topic = "no_op_prefix." + entityTopicName + ".events.v1" // Simulate a prefix
	}

	log.Printf("Kafka (NoOp): Publishing event type %s (data: %+v) to constructed topic %s\n", eventType, data, topic)
	return nil
}
func (p *NoOpProducer) Close() { /* no-op */ }

// Event Data Structures
type UserRegisteredEventData struct {
	UserID                string    `json:"userId"`
	Username              string    `json:"username"`
	Email                 string    `json:"email"`
	Status                string    `json:"status"` // e.g., "pending_verification"
	RegistrationTimestamp time.Time `json:"registrationTimestamp"`
}

type EmailVerificationRequestedEventData struct {
	UserID           string    `json:"userId"`
	Email            string    `json:"email"`
	VerificationCode string    `json:"verificationCode"` // The raw code
	ExpiresAt        time.Time `json:"expiresAt"`
}

type UserEmailVerifiedEventData struct {
	UserID                string    `json:"userId"`
	Email                 string    `json:"email"`
	VerificationTimestamp time.Time `json:"verificationTimestamp"`
}

type UserLoginSuccessEventData struct {
	UserID         string    `json:"userId"`
	SessionID      string    `json:"sessionId,omitempty"` // JTI of refresh token or access token
	IPAddress      string    `json:"ipAddress,omitempty"`
	UserAgent      string    `json:"userAgent,omitempty"`
	LoginTimestamp time.Time `json:"loginTimestamp"`
}

type SessionRevokedEventData struct {
	UserID           string    `json:"userId"`
	SessionID        string    `json:"sessionId,omitempty"` // JTI of the (access) token that was blacklisted/revoked
	RevocationReason string    `json:"revocationReason"`    // "user_logout", "admin_action", "token_compromised"
	RevokedAt        time.Time `json:"revokedAt"`
}

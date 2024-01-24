package auth

import "time"

const (
	HEADERS_CONSUMER_ID           = "X-Consumer-ID"
	HEADERS_CONSUMER_CUSTOM_ID    = "X-Consumer-Custom-ID"
	HEADERS_CONSUMER_USERNAME     = "X-Consumer-Username"
	HEADERS_CREDENTIAL_IDENTIFIER = "X-Credentials-Identifier"
	HEADERS_ANONYMOUS             = "X-Anonymous"
)

type ApiKey struct {
	Key              string
	CreatedAt        time.Time
	TTL              int64
	ConsumerId       int64
	ConsumerCustomId string
	ConsumerUsername string
}

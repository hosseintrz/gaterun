package admin

import "time"

type Consumer struct {
	ID        int64
	Username  string
	CustomID  int64
	CreatedAt time.Time
}

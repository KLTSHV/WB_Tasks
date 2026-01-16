package calendar

import (
	"time"
)

type Event struct {
	ID     int64
	UserID int64
	Date   time.Time
	Text   string
}

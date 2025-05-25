package consumer

import "time"

type Metadata struct {
	Offset     int64
	TimePosted time.Time
}

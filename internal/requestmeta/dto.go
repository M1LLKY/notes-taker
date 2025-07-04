package requestmeta

import "time"

type RequestDTO struct {
	Method    string
	URL       string
	StartTime time.Time
}

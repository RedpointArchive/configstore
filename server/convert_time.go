package main

import (
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

func convertTimeToTimestamp(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

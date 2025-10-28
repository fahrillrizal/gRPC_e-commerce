package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertTimeToTimestamp converts time.Time to protobuf Timestamp
func ConvertTimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

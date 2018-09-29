package hptypes

import (
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

func Timestamp(tsp *timestamp.Timestamp) time.Time {
	if t, err := ptypes.Timestamp(tsp); err == nil {
		return t
	}

	return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
}

func TimestampNow() *timestamp.Timestamp {
	return TimestampProto(time.Now())
}

func TimestampProto(t time.Time) *timestamp.Timestamp {
	if t, err := ptypes.TimestampProto(t); err == nil {
		return t
	}

	t = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))
	return &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
}

package tael

import (
	"errors"
	es "github.com/pingles/tael/elasticsearch"
	"time"
)

type LogMessage struct {
	Message   string
	Timestamp time.Time
	LevelName string
	LogName   string
}

func NewLogMessageFromHit(hit *es.Hit) (*LogMessage, error) {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", hit.Source["@timestamp"].(string))
	if err != nil {
		return nil, err
	}

	message, ok := hit.Source["message"]
	if !ok {
		return nil, errors.New("no message")
	}
	log, ok := hit.Source["logname"]
	if !ok {
		return nil, errors.New("no logname")
	}
	level, ok := hit.Source["level_name"]
	if !ok {
		return nil, errors.New("no level")
	}

	return &LogMessage{
		Message:   message.(string),
		Timestamp: t,
		LogName:   log.(string),
		LevelName: level.(string),
	}, nil
}

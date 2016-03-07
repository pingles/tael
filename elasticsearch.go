package main

import (
	"encoding/json"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/willf/bloom"
	"strings"
	"time"
)

type LogEntry struct {
	Id            string
	Time          time.Time `json:"@timestamp"`
	Message       string    `json:"message"`
	LevelName     string    `json:"level_name"`
	Level         int       `json:"level"`
	ImageName     string    `json:"image_name"`
	ContainerName string    `json:"container_name"`
}

func (e *LogEntry) Header() string {
	return fmt.Sprintf("%s: %-34s [%-5s]", e.Id, e.Time, strings.ToUpper(e.LevelName))
}

func (e *LogEntry) String() string {
	return fmt.Sprintf("%s %s", e.Header(), e.Message)
}

func createSearch(index, query string, now time.Time) *elastigo.SearchDsl {
	from := now.Add(-10 * time.Minute)
	search := elastigo.Search(index)
	fromFilter := elastigo.Filter().Range("@timestamp", from, nil, now, nil, "UTC")
	search.Filter(fromFilter)
	search.Search(query)
	search.Sort(elastigo.Sort("@timestamp").Desc())
	return search
}

func PerformSearch(index, query string, host string) <-chan *LogEntry {
	c := elastigo.NewConn()
	c.SetFromUrl(host)

	filter := bloom.New(20*1000, 5)

	entriesCh := make(chan *LogEntry)
	go func() {
		for {
			search := createSearch(index, query, time.Now())
			result, err := search.Result(c)
			if err != nil {
				return
			}

			for i := len(result.Hits.Hits) - 1; i > 0; i-- {
				hit := result.Hits.Hits[i]
				if filter.Test([]byte(hit.Id)) {
					continue
				}

				var e LogEntry
				err := json.Unmarshal(*hit.Source, &e)
				e.Id = hit.Id

				if err != nil {
					return
				}

				entriesCh <- &e
				filter.Add([]byte(hit.Id))
			}
		}
	}()

	return entriesCh
}

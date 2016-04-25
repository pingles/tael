package elasticsearch

import (
	"errors"
	"encoding/json"
	"fmt"
	"github.com/willf/bloom"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type LogEntry struct {
	Id            string    `json:"id"`
	Time          time.Time `json:"@timestamp"`
	Message       string    `json:"message"`
	LevelName     string    `json:"level_name"`
	Level         int       `json:"level"`
	ImageName     string    `json:"image_name"`
	ContainerName string    `json:"container_name"`
	LogName       string    `json:"logname"`
}

func (e *LogEntry) Header() string {
	return fmt.Sprintf("%s: %-34s [%-5s]", e.Id, e.Time, strings.ToUpper(e.LevelName))
}

func (e *LogEntry) String() string {
	return fmt.Sprintf("%s %s", e.Header(), e.Message)
}



type Hit struct {
	Index  string   `json:"_index"`
	Type   string   `json:"_type"`
	Id     string   `json:"_id"`
	Score  float32  `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type searchHits struct {
	Hits []*Hit `json:"hits"`
}

type searchResp struct {
	Hits *searchHits `json:"hits"`
}

func ExecuteSearch(search *Search) ([]*Hit, error) {
	body, err := search.body()
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(search.url(), "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 200 {
		return nil, errors.New(fmt.Sprintf("unexpected response: %s", string(bytes)))
	}

	var s searchResp
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s.Hits.Hits, nil
}

func StreamSearch(search *Search) <-chan *Hit {
	filter := bloom.New(20*1000, 5)

	entriesCh := make(chan *Hit)
	go func() {
		for _ = range time.Tick(time.Second) {
			results, err := ExecuteSearch(search)
			if err != nil {
				panic(err)
			}

			for i := len(results) - 1; i > 0; i-- {
				hit := results[i]
				idBytes := []byte(hit.Id)

				if filter.Test(idBytes) {
					continue
				}
				entriesCh <- hit
				filter.Add(idBytes)
			}
		}
	}()

	return entriesCh
}

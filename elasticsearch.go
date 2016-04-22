package tael

import (
	"bytes"
	"errors"
	"encoding/json"
	"path"
	"fmt"
	"github.com/willf/bloom"
	"io"
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



type queryString struct {
	Query string `json:"query"`
}

type searchQuery struct {
	QueryString *queryString `json:"query_string"`
}


type searchRequest struct {
	Query *searchQuery           `json:"query"`
	Size  int32                  `json:"size"`
	Sort  map[string]interface{} `json:"sort"`
}


type Search struct {
	Host        string
	Index       string

	// query will be parsed using Lucene syntax
	QueryString string

	Size  int32
	From  int32
}

func (s *Search) url() string {
	return fmt.Sprintf("%s/%s", s.Host, path.Join(s.Index, "_search"))
}

// builds the search query payload
func (s *Search) query() *searchRequest {
	return &searchRequest{
		Query: &searchQuery{
			QueryString: &queryString{
				Query: s.QueryString,
			},
		},
		Size: s.Size,
		Sort: map[string]interface{}{
			"@timestamp": map[string]string {
				"order": "desc",
				"unmapped_type": "boolean",
			},
		},
	}
}

func (s *Search) body() (io.Reader, error) {
	bs, err := json.Marshal(s.query())
	if err != nil {
		return nil, err
	}
	fmt.Println("QUERY:", string(bs))
	return bytes.NewReader(bs), nil
}

const (
	DefaultSize = 20
)

func NewSearchWithFilters(host, index, query string, now time.Time, filters map[string]string) *Search {
	return &Search{
		Host: host,
		Index: index,
		Size: DefaultSize,
		QueryString: query,
	}
}

func NewSearch(host, index, query string, now time.Time) *Search {
	return &Search{
		Host: host,
		Index: index,
		Size: DefaultSize,
		QueryString: query,
	}
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
			fmt.Println("running search")
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

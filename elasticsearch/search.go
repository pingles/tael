package elasticsearch

import (
	"fmt"
	"path"
	"io"
	"encoding/json"
	"bytes"
	"time"
)


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

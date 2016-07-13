package store

import (
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"net/url"
	"strings"
)

// NOTICE: do not forget to setup elastic search client

// ESClient elastic search client
var ESClient *elastic.Client
var ElasticSearchIndexName = "fader"

// --------------------
// SearchService
// --------------------

type SearchService struct {
	searchSource interface{}
	client       *elastic.Client
	index        []string
	typ          []string
}

func MustSearchService() *SearchService {
	return NewSearchService(ESClient)
}

// NewSearchService
func NewSearchService(client *elastic.Client) *SearchService {

	return &SearchService{
		client: client,
	}
}

func (s *SearchService) SetQueryRaw(src interface{}) *SearchService {
	s.searchSource = src
	return s
}

// Index
func (s *SearchService) Index(index ...string) *SearchService {
	if s.index == nil {
		s.index = make([]string, 0)
	}
	s.index = append(s.index, index...)
	return s
}

// Types
func (s *SearchService) Type(typ ...string) *SearchService {
	if s.typ == nil {
		s.typ = make([]string, 0)
	}
	s.typ = append(s.typ, typ...)
	return s
}

// Do executes the search and returns a SearchResult.
func (s *SearchService) Do() (*elastic.SearchResult, error) {
	path := fmt.Sprintf("/%s/%s/_search", strings.Join(s.index, ","), strings.Join(s.typ, ","))

	res, err := s.client.PerformRequest("POST", path, url.Values{}, s.searchSource)
	if err != nil {
		return nil, err
	}

	// Return search results
	ret := new(elastic.SearchResult)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

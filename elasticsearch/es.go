package elasticsearch

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type Client struct {
	url string
}

func (c *Client) Index(name string) *Index {

	return &Index{
		name:   name,
		client: c,
	}
}

func (c *Client) Search(index string, dsl interface{}) *Result {
	request := gorequest.New()
	url := fmt.Sprintf("%s%s/_search", c.url, index)

	resultBody := &ResultBody{}

	_, _, errs := request.Post(url).Send(&dsl).EndStruct(&resultBody)

	return &Result{
		body:   resultBody,
		errors: errs,
	}
}

type Index struct {
	client *Client
	name   string
}

type Aggregation struct {
	Buckets      []interface{}
	Aggregations interface{}
}

type ResultBody struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int     `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string                 `json:"_index"`
			Type   string                 `json:"_type"`
			ID     string                 `json:"_id"`
			Score  float64                `json:"_score"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations interface{} `json:"aggregations"`
}

type Result struct {
	body   *ResultBody
	errors []error
}

func (r *Result) Body() *ResultBody {
	return r.body
}

func (r *Result) Errors() []error {
	return r.errors
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

package elasticsearch

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type Client struct {
	url string
}

func (c *Client) Search(index string, dsl interface{}) *Result {
	return c.post(index, "search", dsl)
}

func (c *Client) Count(index string, dsl interface{}) *Result {
	return c.post(index, "count", dsl)
}

func (c *Client) post(index string, action string, dsl interface{}) *Result {
	request := gorequest.New()
	url := fmt.Sprintf("%s%s/_%s", c.url, index, action)

	resultBody := &ResultBody{}

	resp, body, errs := request.Post(url).Send(&dsl).EndStruct(&resultBody)

	if resp.StatusCode >= 400 {
		if errs == nil {
			errs = []error{}
		}
		errs = append(errs, errors.New(string(body)))
	}

	return &Result{
		body:   resultBody,
		errors: errs,
	}
}

type Aggregation struct {
	Buckets      []interface{}
	Aggregations interface{}
}

type ResultBody struct {
	Count    int  `json:"count"`
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

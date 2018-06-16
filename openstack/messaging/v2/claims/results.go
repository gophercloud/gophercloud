package claims

import "github.com/gophercloud/gophercloud"

func (r CreateResult) Extract() ([]Messages, error) {
	var s struct {
		Messages []Messages `json:"messages"`
	}
	err := r.ExtractInto(&s)
	return s.Messages, err
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	gophercloud.Result
}

type Messages struct {
	Age  float32                `json:"age"`
	Href string                 `json:"href"`
	TTL  int                    `json:"ttl"`
	Body map[string]interface{} `json:"body"`
}

package queues

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// QueuePage contains a single page of all queues from a List operation.
type QueuePage struct {
	pagination.LinkedPageBase
}

// CreateResult is the response of a Create operation.
type CreateResult struct {
	gophercloud.ErrResult
}

// UpdateResult is the response of a Update operation.
type UpdateResult struct {
	commonResult
}

// GetResult is the response of a Get operation.
type GetResult struct {
	commonResult
}

// Queue represents a messaging queue.
type Queue struct {
	Href          string       `json:"href"`
	Methods       []string     `json:"methods"`
	Name          string       `json:"name"`
	Paths         []string     `json:"paths"`
	ResourceTypes []string     `json:"resource_types"`
	Metadata      QueueDetails `json:"metadata"`
}

// QueueDetails represents the metadata of a queue.
type QueueDetails struct {
	DeadLetterQueue           string `json:"_dead_letter_queue"`
	DeadLetterQueueMessageTTL int    `json:"_dead_letter_queue_messages_ttl"`
	DefaultMessageDelay       int    `json:"_default_message_delay"`
	DefaultMessageTTL         int    `json:"_default_message_ttl"`
	MaxClaimCount             int    `json:"_max_claim_count"`
	MaxMessagesPostSize       int    `json:"_max_messages_post_size"`
	Description               string `json:"description"`
	Flavor                    string `json:"flavor"`
}

// Extract interprets any commonResult as a Queue.
func (r commonResult) Extract() (QueueDetails, error) {
	var s QueueDetails
	err := r.ExtractInto(&s)
	return s, err
}

// ExtractQueues interprets the results of a single page from a
// List() call, producing a map of queues.
func ExtractQueues(r pagination.Page) ([]Queue, error) {
	var s struct {
		Queues []Queue `json:"queues"`
	}
	err := (r.(QueuePage)).ExtractInto(&s)
	return s.Queues, err
}

// IsEmpty determines if a QueuesPage contains any results.
func (page QueuePage) IsEmpty() (bool, error) {
	s, err := ExtractQueues(page)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (page QueuePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

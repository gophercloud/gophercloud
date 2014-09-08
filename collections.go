package gophercloud

import (
	"errors"

	"github.com/racker/perigee"
)

var (
	// ErrPageNotAvailable is returned from a Pager when a next or previous page is requested, but does not exist.
	ErrPageNotAvailable = errors.New("The requested Collection page does not exist.")
)

// Collection describes the minimum functionality that any collection resource must implement to be able to use
// the global paging and iteration functions.
// Every resource that returns a list of multiple results must implement this functionality, whether or not it is paged.
// In addition to the methods provided here, each collection should also provide an AsItem(Page) method that
// casts the Page to its more specific type and returns the Page's contents as a slice.
type Collection interface {

	// Pager returns one of the concrete Pager implementations from this package, or a custom one.
	// The style of Pager returned determines how the collection is paged.
	Pager() Pager
}

// EachPage iterates through a Collection one page at a time.
// The handler function will be invoked with a Collection containing each page.
// If the handler returns true, iteration will continue. If it returns false, no more pages will be fetched.
func EachPage(first Collection, handler func(Collection) bool) error {
	p := first.Pager()
	var err error
	current := first

	for {
		if !handler(current) {
			return nil
		}

		if !p.HasNextPage() {
			return nil
		}

		current, err = p.NextPage()
		if err != nil {
			return err
		}
	}
}

// AllPages consolidates all pages reachable from a provided starting point into a single mega-Page.
// Use this only when you know that the full set will always fit within memory.
func AllPages(first Collection) (Collection, error) {
	return first, nil
}

// Pager describes a specific paging idiom for a Page resource.
// Generally, to use a Pager, the Page must also implement a more specialized interface than Page.
// Clients should not generally interact with Pagers directly.
// Instead, use the more convenient collection traversal methods: All and EachPage.
type Pager interface {

	// HasNextPage returns true if a call to NextPage will return an additional Page of results.
	HasNextPage() bool

	// NextPage returns the next Page in the sequence.
	// Panics if no page is available, so always check HasNextPage first.
	NextPage() (Collection, error)
}

// SinglePager is used by collections that are not actually paged.
// It has no additional interface requirements for its host Page.
type SinglePager struct{}

// HasNextPage always reports false.
func (p SinglePager) HasNextPage() bool {
	return false
}

// NextPage always returns an ErrPageNotAvailable.
func (p SinglePager) NextPage() (Collection, error) {
	return nil, ErrPageNotAvailable
}

// PaginationLinks stores the `next` and `previous` links that are provided by some (but not all) paginated resources.
type PaginationLinks struct {

	// Next is the full URL to the next page of results, or nil if this is the last page.
	Next *string `json:"next,omitempty"`

	// Previous is the full URL to the previous page of results, or nil if this is the first page.
	Previous *string `json:"previous,omitempty"`
}

// LinkCollection must be satisfied by a Page that uses a LinkPager.
type LinkCollection interface {
	Collection

	// Service returns the client used to make further requests.
	Service() *ServiceClient

	// Links returns the pagination links from a single page.
	Links() PaginationLinks

	// Interpret an arbitrary JSON result as a new LinkCollection.
	Interpret(interface{}) (LinkCollection, error)
}

// LinkPager implements paging for collections that provide a link structure in their response JSON.
// It follows explicit `next` links and stops when the `next` link is "null".
type LinkPager struct {
	current LinkCollection
}

// NewLinkPager creates and initializes a pager for a LinkCollection.
func NewLinkPager(first LinkCollection) *LinkPager {
	return &LinkPager{current: first}
}

// HasNextPage checks the `next` link in the pagination data.
func (p *LinkPager) HasNextPage() bool {
	return p.current.Links().Next != nil
}

// NextPage follows the `next` link to construct the next page of data.
func (p *LinkPager) NextPage() (Collection, error) {
	url := p.current.Links().Next
	if url == nil {
		return nil, ErrPageNotAvailable
	}

	var response interface{}
	_, err := perigee.Request("GET", *url, perigee.Options{
		MoreHeaders: p.current.Service().Provider.AuthenticatedHeaders(),
		Results:     &response,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	interpreted, err := p.current.Interpret(response)
	if err != nil {
		return nil, err
	}

	p.current = interpreted
	return interpreted, nil
}

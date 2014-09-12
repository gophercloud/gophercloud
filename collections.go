package gophercloud

import "errors"

var (
	// ErrPageNotAvailable is returned from a Pager when a next or previous page is requested, but does not exist.
	ErrPageNotAvailable = errors.New("The requested Collection page does not exist.")
)

// Collection must be satisfied by the result type of any resource collection.
// It allows clients to interact with the resource uniformly, regardless of whether or not or how it's paginated.
type Collection interface {

	// NextPageURL generates the URL for the page of data that follows this collection.
	// Return "" if no such page exists.
	NextPageURL() string

	// Concat creates a new Collection that contains all of the elements from this page and another page.
	// It's used to aggregate results for the AllPages method.
	Concat(Collection) Collection
}

// Pager knows how to advance through a specific resource collection, one page at a time.
type Pager struct {
	initialURL string

	advance func(string) (Collection, error)
}

// NewPager constructs a manually-configured pager.
// Supply the URL for the first page and a function that requests a specific page given a URL.
func NewPager(initialURL string, advance func(string) (Collection, error)) Pager {
	return Pager{
		initialURL: initialURL,
		advance:    advance,
	}
}

// NewSinglePager constructs a Pager that "iterates" over a single-paged Collection.
// Supply a function that returns the only page.
func NewSinglePager(only func() (Collection, error)) Pager {
	consumed := false
	single := func(_ string) (Collection, error) {
		if !consumed {
			consumed = true
			return only()
		}
		return nil, ErrPageNotAvailable
	}

	return Pager{
		initialURL: "",
		advance:    single,
	}
}

// EachPage iterates over each page returned by a Pager, yielding one at a time to a handler function.
// Return "false" from the handler to prematurely stop iterating.
func (p Pager) EachPage(handler func(Collection) bool) error {
	currentURL := p.initialURL
	for {
		currentPage, err := p.advance(currentURL)
		if err != nil {
			return err
		}

		if !handler(currentPage) {
			return nil
		}

		currentURL = currentPage.NextPageURL()
		if currentURL == "" {
			return nil
		}
	}
}

// AllPages accumulates every page reachable from a Pager into a single Collection, for convenience.
func (p Pager) AllPages() (Collection, error) {
	var megaPage Collection

	err := p.EachPage(func(page Collection) bool {
		if megaPage == nil {
			megaPage = page
		} else {
			megaPage = megaPage.Concat(page)
		}
		return true
	})

	return megaPage, err
}

// PaginationLinks stores the `next` and `previous` links that are provided by some (but not all) paginated resources.
type PaginationLinks struct {

	// Next is the full URL to the next page of results, or nil if this is the last page.
	Next *string `json:"next,omitempty"`

	// Previous is the full URL to the previous page of results, or nil if this is the first page.
	Previous *string `json:"previous,omitempty"`
}

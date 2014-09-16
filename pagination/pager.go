package pagination

import "errors"

var (
	// ErrPageNotAvailable is returned from a Pager when a next or previous page is requested, but does not exist.
	ErrPageNotAvailable = errors.New("The requested page does not exist.")
)

// Page must be satisfied by the result type of any resource collection.
// It allows clients to interact with the resource uniformly, regardless of whether or not or how it's paginated.
// Generally, rather than implementing this interface directly, implementors should embed one of the concrete PageBase structs,
// instead.
// Depending on the pagination strategy of a particular resource, there may be an additional subinterface that the result type
// will need to implement.
type Page interface {

	// NextPageURL generates the URL for the page of data that follows this collection.
	// Return "" if no such page exists.
	NextPageURL() (string, error)

	// IsEmpty returns true if this Page has no items in it.
	IsEmpty() (bool, error)
}

// Pager knows how to advance through a specific resource collection, one page at a time.
type Pager struct {
	initialURL string

	fetchNextPage func(string) (Page, error)
}

// NewPager constructs a manually-configured pager.
// Supply the URL for the first page, a function that requests a specific page given a URL, and a function that counts a page.
func NewPager(initialURL string, fetchNextPage func(string) (Page, error)) Pager {
	return Pager{
		initialURL:    initialURL,
		fetchNextPage: fetchNextPage,
	}
}

// EachPage iterates over each page returned by a Pager, yielding one at a time to a handler function.
// Return "false" from the handler to prematurely stop iterating.
func (p Pager) EachPage(handler func(Page) (bool, error)) error {
	currentURL := p.initialURL
	for {
		currentPage, err := p.fetchNextPage(currentURL)
		if err != nil {
			return err
		}

		empty, err := currentPage.IsEmpty()
		if err != nil {
			return err
		}
		if empty {
			return nil
		}

		ok, err := handler(currentPage)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		currentURL, err = currentPage.NextPageURL()
		if err != nil {
			return err
		}
		if currentURL == "" {
			return nil
		}
	}
}

package pager

import (
	"fmt"
	"maps"
	"math"

	"crdx.org/lighthouse/pkg/util/webutil"
)

const Key = "p"

type State struct {
	CurrentPage     uint
	TotalPages      uint
	FirstPageURL    string
	NextPageURL     string
	PreviousPageURL string
	LastPageURL     string
}

// GetCurrentPageNumber returns the page number of the current page, and true if a valid page number
// was provided.
func GetCurrentPageNumber(pageNumber uint) (uint, bool) {
	if pageNumber < 1 {
		return 0, false
	}

	return pageNumber, true
}

// GetState returns an instance of *State with paging information the template needs to render the
// paging navigation (p/pager/nav.go.html), or an error if basePath could not be parsed.
//
// This method will not modify qs. Pass nil for qs if no additional parameters are needed.
func GetState(pageNumber uint, pageCount uint, basePath string, qs map[string]string) (*State, error) {
	if qs == nil {
		qs = map[string]string{}
	} else {
		qs = maps.Clone(qs)
	}

	state := &State{
		CurrentPage: pageNumber,
		TotalPages:  pageCount,
	}

	var err error

	// Removes a lot of the repetitive error handling from the code below using a technique inspired
	// by https://go.dev/blog/errors-are-values.
	f := func(n uint) string {
		if err != nil {
			return ""
		}

		qs[Key] = fmt.Sprint(n)

		var url string
		if url, err = webutil.BuildURL(basePath, qs); err != nil {
			return ""
		}

		return url
	}

	if pageNumber < pageCount {
		state.NextPageURL = f(pageNumber + 1)
		state.LastPageURL = f(pageCount)
	}

	if pageNumber > 1 {
		state.PreviousPageURL = f(pageNumber - 1)
		state.FirstPageURL = f(1)
	}

	return state, err
}

// GetPageCount returns the number of pages needed to fit n items if there are perPage items per
// page. If there are no items then that's still 1 page needed to show nothing.
func GetPageCount(n uint, perPage uint) uint {
	if n == 0 {
		return 1
	}
	return uint(math.Ceil(float64(n) / float64(perPage)))
}

// GetOffset returns the offset for a LIMIT query for page pageNumber if there are perPage items per
// page.
func GetOffset(pageNumber uint, perPage uint) uint {
	return (pageNumber - 1) * perPage
}

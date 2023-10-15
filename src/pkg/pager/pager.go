package pager

import (
	"fmt"
	"maps"
	"math"

	"crdx.org/lighthouse/util/webutil"
	"github.com/gofiber/fiber/v2"
)

const queryStringParameter = "page"

type State struct {
	CurrentPage     uint
	TotalPages      uint
	NextPageURL     string
	PreviousPageURL string
}

// GetCurrentPageNumber returns the page number of the current page, and true if a valid page number
// was provided.
func GetCurrentPageNumber(c *fiber.Ctx) (uint, bool) {
	pageNumber := uint(c.QueryInt(queryStringParameter, 1))

	if pageNumber < 1 {
		return 0, false
	}

	return pageNumber, true
}

// GetState returns an instance of *State with paging information the template needs to render the
// paging navigation (p/pager/nav.go.html), or an error if basePath could not be parsed.
//
// This method will not modify qs.
func GetState(pageNumber uint, pageCount uint, basePath string, qs map[string]string) (*State, error) {
	qs = maps.Clone(qs)

	state := &State{
		CurrentPage: pageNumber,
		TotalPages:  pageCount,
	}

	var err error

	if pageNumber < pageCount {
		qs[queryStringParameter] = fmt.Sprint(pageNumber + 1)
		state.NextPageURL, err = webutil.BuildURL(basePath, qs)
		if err != nil {
			return nil, err
		}
	}

	if pageNumber > 1 {
		qs[queryStringParameter] = fmt.Sprint(pageNumber - 1)
		state.PreviousPageURL, err = webutil.BuildURL(basePath, qs)
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

// GetPageCount returns the number of pages needed to fit n items if there are perPage items per
// page.
func GetPageCount(n uint, perPage uint) uint {
	return uint(math.Ceil(float64(n) / float64(perPage)))
}

// GetOffset returns the offset for a LIMIT query for page pageNumber if there are perPage items per
// page.
func GetOffset(pageNumber uint, perPage uint) uint {
	return (pageNumber - 1) * perPage
}

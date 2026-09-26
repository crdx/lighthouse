package fontawesome

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

type results struct {
	Data struct {
		Search []icon `json:"search"`
	} `json:"data"`
}

type icon struct {
	FamilyStylesByLicense familyStylesByLicense `json:"familyStylesByLicense"`
	ID                    string                `json:"id"`
	Label                 string                `json:"label"`
}

type familyStylesByLicense struct {
	Pro []familyStyle `json:"pro"`
}

type familyStyle struct {
	Family string `json:"family"`
	Style  string `json:"style"`
}

const (
	styleDuotone = "duotone"
	styleSolid   = "solid"
	styleBrands  = "brands"
)

var availableStyles = []string{styleDuotone, styleSolid, styleBrands}

const (
	maxResults    = 10
	upstreamLimit = 50
)

var apiURL = "https://api.fontawesome.com"

// Search searches the FontAwesome API for an icon matching the query.
func Search(query string) ([]map[string]string, bool, error) {
	var wantedStyle string

	if strings.Contains(query, ":") {
		wantedStyle, query, _ = strings.Cut(query, ":")

		if !slices.Contains(availableStyles, wantedStyle) {
			return []map[string]string{}, false, nil
		}
	}

	found, err := search(query)
	if err != nil {
		return nil, false, err
	}

	want := func(style string) bool {
		return wantedStyle == "" || wantedStyle == style
	}

	icons := []map[string]string{}

	for _, icon := range found {
		add := func(style string) {
			icons = append(icons, map[string]string{
				"style": style,
				"name":  icon.ID,
				"label": icon.Label,
			})
		}

		for _, style := range icon.FamilyStylesByLicense.Pro {
			if isDuotone(style) && want(styleDuotone) {
				add(styleDuotone)
			}

			if isSolid(style) && want(styleSolid) {
				add(styleSolid)
			}

			if isBrands(style) && want(styleBrands) {
				add(styleBrands)
			}
		}
	}

	if len(icons) > maxResults {
		return icons[:maxResults], true, nil
	}

	return icons, len(found) == upstreamLimit, nil
}

func search(s string) ([]icon, error) {
	if s == "" {
		return []icon{}, nil
	}

	q := `
		query {
			search(
				version: "6.4.2",
				query: %s,
				first: %d
			) {
				id,
				label,
				familyStylesByLicense {
					pro {
						family,
						style
					},
				}
			}
		}
	`

	payload := fmt.Sprintf(q, strconv.Quote(s), upstreamLimit)
	response, err := req.R().
		SetBodyJsonMarshal(map[string]string{"query": payload}).
		Post(apiURL)
	if err != nil {
		return nil, fmt.Errorf("fontawesome request: %w", err)
	}

	if !response.IsSuccessState() {
		return nil, fmt.Errorf("fontawesome request: unexpected status %d", response.StatusCode)
	}

	var results results
	if err := json.Unmarshal(response.Bytes(), &results); err != nil {
		return nil, fmt.Errorf("fontawesome response: %w", err)
	}

	return results.Data.Search, nil
}

func isSolid(style familyStyle) bool {
	return style.Family == "classic" && style.Style == "solid"
}

func isDuotone(style familyStyle) bool {
	return style.Family == "duotone"
}

func isBrands(style familyStyle) bool {
	return style.Family == "classic" && style.Style == "brands"
}

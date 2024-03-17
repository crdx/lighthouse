package fontawesome

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
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

const maxResults = 10

// Search searches the FontAwesome API for an icon matching the query.
func Search(query string) ([]map[string]string, bool) {
	var wantedStyle string

	if strings.Contains(query, ":") {
		wantedStyle, query, _ = strings.Cut(query, ":")

		if !slices.Contains(availableStyles, wantedStyle) {
			return []map[string]string{}, false
		}
	}

	want := func(style string) bool {
		return wantedStyle == "" || wantedStyle == style
	}

	icons := []map[string]string{}

	for _, icon := range search(query) {
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
		return icons[:maxResults], true
	}

	return icons, false
}

func search(s string) []icon {
	if s == "" {
		return []icon{}
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

	payload := fmt.Sprintf(q, strconv.Quote(s), maxResults)
	jsonBytes := lo.Must(json.Marshal(map[string]string{"query": payload}))
	response := lo.Must(req.R().SetBodyJsonBytes(jsonBytes).Post("https://api.fontawesome.com"))

	var results results
	lo.Must0(json.Unmarshal(response.Bytes(), &results))

	return results.Data.Search
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

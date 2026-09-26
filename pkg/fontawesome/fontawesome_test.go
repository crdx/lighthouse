package fontawesome

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	classicSolid  = familyStyle{Family: "classic", Style: "solid"}
	classicLight  = familyStyle{Family: "classic", Style: "light"}
	classicBrands = familyStyle{Family: "classic", Style: "brands"}
	duotoneSolid  = familyStyle{Family: "duotone", Style: "solid"}
)

func makeIcons(count int, styles ...familyStyle) []icon {
	icons := make([]icon, count)
	for i := range icons {
		icons[i] = icon{
			ID:                    fmt.Sprintf("icon-%d", i),
			Label:                 fmt.Sprintf("Icon %d", i),
			FamilyStylesByLicense: familyStylesByLicense{Pro: styles},
		}
	}
	return icons
}

func serve(t *testing.T, icons []icon) *string {
	t.Helper()

	var requestBody string

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		requestBody = string(body)

		var response results
		response.Data.Search = icons
		_ = json.NewEncoder(writer).Encode(response)
	}))

	setAPIURL(t, server.URL)
	t.Cleanup(server.Close)

	return &requestBody
}

func setAPIURL(t *testing.T, url string) {
	t.Helper()

	original := apiURL
	apiURL = url
	t.Cleanup(func() { apiURL = original })
}

func names(icons []map[string]string) []string {
	result := make([]string, len(icons))
	for i, icon := range icons {
		result[i] = icon["style"] + ":" + icon["name"]
	}
	return result
}

func TestSearchExpandsStyles(t *testing.T) {
	requestBody := serve(t, makeIcons(2, classicSolid, classicLight, duotoneSolid))

	icons, hasMore, err := Search("icon")

	require.NoError(t, err)
	assert.Equal(t, []string{"solid:icon-0", "duotone:icon-0", "solid:icon-1", "duotone:icon-1"}, names(icons))
	assert.False(t, hasMore)
	assert.Contains(t, *requestBody, fmt.Sprintf("first: %d", upstreamLimit))
	assert.Contains(t, *requestBody, `query: \"icon\"`)
}

func TestSearchFiltersByStyle(t *testing.T) {
	icons := append(makeIcons(3, classicSolid, duotoneSolid), icon{
		ID:                    "github",
		Label:                 "GitHub",
		FamilyStylesByLicense: familyStylesByLicense{Pro: []familyStyle{classicBrands}},
	})
	requestBody := serve(t, icons)

	found, hasMore, err := Search("brands:git")

	require.NoError(t, err)
	assert.Equal(t, []string{"brands:github"}, names(found))
	assert.False(t, hasMore)
	assert.Contains(t, *requestBody, `query: \"git\"`)
}

func TestSearchTruncatesToMaxResults(t *testing.T) {
	serve(t, makeIcons(maxResults, classicSolid, duotoneSolid))

	icons, hasMore, err := Search("icon")

	require.NoError(t, err)
	assert.Len(t, icons, maxResults)
	assert.True(t, hasMore)
}

func TestSearchFullUpstreamPageHasMore(t *testing.T) {
	icons := makeIcons(upstreamLimit, classicSolid)
	icons[0].FamilyStylesByLicense.Pro = []familyStyle{duotoneSolid}
	serve(t, icons)

	found, hasMore, err := Search("duotone:icon")

	require.NoError(t, err)
	assert.Equal(t, []string{"duotone:icon-0"}, names(found))
	assert.True(t, hasMore)
}

func TestSearchSkipsRequest(t *testing.T) {
	for _, query := range []string{"", "light:icon", "duotone:"} {
		t.Run(query, func(t *testing.T) {
			requestBody := serve(t, makeIcons(1, classicSolid))

			icons, hasMore, err := Search(query)

			require.NoError(t, err)
			assert.Empty(t, icons)
			assert.False(t, hasMore)
			assert.Empty(t, *requestBody)
		})
	}
}

func TestSearchUpstreamErrors(t *testing.T) {
	testCases := map[string]http.HandlerFunc{
		"status": func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		},
		"json": func(writer http.ResponseWriter, _ *http.Request) {
			_, _ = io.Copy(writer, strings.NewReader("<html>"))
		},
	}

	for name, handler := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(handler)
			t.Cleanup(server.Close)
			setAPIURL(t, server.URL)

			icons, hasMore, err := Search("icon")

			require.Error(t, err)
			assert.Empty(t, icons)
			assert.False(t, hasMore)
		})
	}
}

func TestSearchUnreachable(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	server.Close()
	setAPIURL(t, server.URL)

	_, _, err := Search("icon")

	require.ErrorContains(t, err, "fontawesome request")
}

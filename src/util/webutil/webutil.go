package webutil

import "net/url"

// BuildURL builds a URL out of basePath and the contents of queryParams, or an error if basePath
// could not be parsed.
func BuildURL(basePath string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return "", err
	}

	values := url.Values{}

	for key, value := range queryParams {
		values.Add(key, value)
	}

	u.RawQuery = values.Encode()
	return u.String(), nil
}

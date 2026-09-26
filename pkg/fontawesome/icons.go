package fontawesome

import (
	"regexp"
	"slices"
	"strings"
)

var (
	ruleRegexp     = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)
	selectorRegexp = regexp.MustCompile(`^\.fa-([a-z0-9-]+)::?before$`)
)

var iconNames = map[string]bool{}

func Init(stylesheet []byte) {
	iconNames = parseIconNames(stylesheet)
}

func parseIconNames(stylesheet []byte) map[string]bool {
	names := map[string]bool{}

	for _, match := range ruleRegexp.FindAllSubmatch(stylesheet, -1) {
		selectors, declarations := string(match[1]), string(match[2])
		if !strings.Contains(declarations, "content:") {
			continue
		}

		for selector := range strings.SplitSeq(selectors, ",") {
			if nameMatch := selectorRegexp.FindStringSubmatch(strings.TrimSpace(selector)); nameMatch != nil {
				names[nameMatch[1]] = true
			}
		}
	}

	return names
}

func IsValidIcon(value string) bool {
	style, name, found := strings.Cut(value, ":")
	return found && slices.Contains(availableStyles, style) && iconNames[name]
}

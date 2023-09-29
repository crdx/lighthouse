package tpl

type SortableColumnConfig struct {
	Label                string
	DefaultSortDirection string
}

type SortableColumnState struct {
	Label                string
	CurrentSortColumn    string
	CurrentSortDirection string
	SortColumn           string
	SortDirection        string
}

// AddSortMetadata takes information about table columns and the current sort column and direction
// and adds metadata for the template to use to render the sortable table headings.
func AddSortMetadata(currentSortColumn string, currentSortDirection string, input map[string]SortableColumnConfig) map[string]SortableColumnState {
	output := map[string]SortableColumnState{}

	for slug, column := range input {
		sortDirection := column.DefaultSortDirection

		if currentSortColumn == slug {
			if currentSortDirection == "desc" {
				sortDirection = "asc"
			} else {
				sortDirection = "desc"
			}
		}

		output[slug] = SortableColumnState{
			Label:                column.Label,
			CurrentSortColumn:    currentSortColumn,
			CurrentSortDirection: currentSortDirection,
			SortColumn:           slug,
			SortDirection:        sortDirection,
		}
	}

	return output
}

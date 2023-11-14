package tplutil

type ColumnConfig struct {
	Label                any
	DefaultSortDirection string
	Minimal              bool
}

type ColumnState struct {
	Label                any
	CurrentFilter        string
	CurrentSortColumn    string
	CurrentSortDirection string
	SortColumn           string
	SortDirection        string
	Minimal              bool
}

// AddMetadata takes information about table columns and the current sort column and direction
// and adds metadata for the template to use to render the sortable table headings.
func AddMetadata(currentSortColumn string, currentSortDirection string, currentFilter string, input map[string]ColumnConfig) map[string]ColumnState {
	output := map[string]ColumnState{}

	for slug, column := range input {
		sortDirection := column.DefaultSortDirection

		if currentSortColumn == slug {
			if currentSortDirection == "desc" {
				sortDirection = "asc"
			} else {
				sortDirection = "desc"
			}
		}

		output[slug] = ColumnState{
			Label:                column.Label,
			CurrentSortColumn:    currentSortColumn,
			CurrentSortDirection: currentSortDirection,
			CurrentFilter:        currentFilter,
			SortColumn:           slug,
			SortDirection:        sortDirection,
			Minimal:              column.Minimal,
		}
	}

	return output
}

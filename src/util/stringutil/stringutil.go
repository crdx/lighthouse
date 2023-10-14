package stringutil

func Pluralise(count int, unit string) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
}

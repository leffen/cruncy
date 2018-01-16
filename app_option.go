package cruncy

import "strings"

// ParseNames splits a input string on , trimsspaces and removes all prefixed with #
// Makes it easier to use list values as (multipline) input params
func ParseNames(input string) []string {
	rc := []string{}
	items := strings.FieldsFunc(input, func(r rune) bool { return r == ',' || r == '\n' })
	for _, item := range items {
		stripped := strings.TrimSpace(item)
		if len(stripped) < 1 {
			continue
		}
		if stripped[0] == '#' {
			continue
		}
		rc = append(rc, stripped)
	}
	return rc
}

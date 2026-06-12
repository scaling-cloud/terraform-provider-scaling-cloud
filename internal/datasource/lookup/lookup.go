// Package lookup resolves a single record out of a listed collection by
// matching one of its fields. Every Scaling data source is a name/email→id
// lookup over a list endpoint, so they share this client-side filter rather
// than each re-deriving the not-found and ambiguous-match diagnostics.
package lookup

import "fmt"

// One returns the single item whose key equals value. kind is the human label
// for the record type (for example "inbound integration") and field is the
// matched attribute (for example "name"); both appear in error messages.
//
// It errors when nothing matches, so a typo'd lookup fails loudly, and when
// more than one item matches, so an ambiguous lookup never silently resolves to
// an arbitrary row.
func One[T any](items []T, kind, field, value string, key func(T) string) (T, error) {
	var match T
	found := false
	for _, item := range items {
		if key(item) != value {
			continue
		}
		if found {
			var zero T
			return zero, fmt.Errorf(
				"multiple %s found with %s %q; it must be unique to look one up",
				kind, field, value,
			)
		}
		match = item
		found = true
	}
	if !found {
		return match, fmt.Errorf("no %s found with %s %q", kind, field, value)
	}
	return match, nil
}

package lookup

import (
	"strings"
	"testing"
)

type item struct {
	id   string
	name string
}

func key(i item) string { return i.name }

func TestOneReturnsSingleMatch(t *testing.T) {
	t.Parallel()

	items := []item{
		{id: "a", name: "alpha"},
		{id: "b", name: "beta"},
	}

	got, err := One(items, "thing", "name", "beta", key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.id != "b" {
		t.Errorf("got id %q, want b", got.id)
	}
}

func TestOneErrorsWhenNoMatch(t *testing.T) {
	t.Parallel()

	_, err := One([]item{{id: "a", name: "alpha"}}, "thing", "name", "missing", key)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	// The message must name the kind, the field, and the value so a typo is
	// obvious from the diagnostic alone.
	for _, want := range []string{"no thing", "name", "missing"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err.Error(), want)
		}
	}
}

func TestOneErrorsWhenAmbiguous(t *testing.T) {
	t.Parallel()

	items := []item{
		{id: "a", name: "dup"},
		{id: "b", name: "dup"},
	}

	_, err := One(items, "thing", "name", "dup", key)
	if err == nil {
		t.Fatal("expected an error for an ambiguous match, got nil")
	}
	for _, want := range []string{"multiple", "thing", "dup"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err.Error(), want)
		}
	}
}

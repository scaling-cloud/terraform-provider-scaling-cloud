package user

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func strPtr(s string) *string { return &s }

func TestMapToModel(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.User{
		ID:        "user_1",
		FirstName: strPtr("Jane"),
		LastName:  strPtr("Smith"),
		Email:     strPtr("jane@example.com"),
	})

	if got.ID.ValueString() != "user_1" {
		t.Errorf("ID = %q, want user_1", got.ID.ValueString())
	}
	if got.Email.ValueString() != "jane@example.com" {
		t.Errorf("Email = %q, want jane@example.com", got.Email.ValueString())
	}
	if got.FirstName.ValueString() != "Jane" {
		t.Errorf("FirstName = %q, want Jane", got.FirstName.ValueString())
	}
}

func TestMapToModelNullNames(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.User{ID: "user_1", Email: strPtr("x@y.z")})
	// Unknown name fields must surface as Terraform nulls, not empty strings.
	if !got.FirstName.IsNull() || !got.LastName.IsNull() {
		t.Errorf("FirstName/LastName = %v/%v, want null", got.FirstName, got.LastName)
	}
}

func TestEmailOf(t *testing.T) {
	t.Parallel()

	if emailOf(client.User{Email: strPtr("a@b.c")}) != "a@b.c" {
		t.Errorf("emailOf with email should return it")
	}
	// A user with no email can never match a requested email.
	if emailOf(client.User{Email: nil}) != "" {
		t.Errorf("emailOf with nil email should return empty string")
	}
}

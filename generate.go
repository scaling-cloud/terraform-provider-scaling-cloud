// Package main hosts the provider entrypoint and documentation generation.
//
// Run `go generate ./...` to regenerate the docs/ directory from the provider
// schema and the examples/ directory. Documentation is published to the
// Terraform Registry, so regenerate and commit it whenever a resource schema
// or its example changes.
//
// tfplugindocs is run via `go run …@version` rather than added as a module
// dependency, so it stays out of go.mod and cannot drag the module's Go
// language version forward.
package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.23.0 generate --provider-name scaling

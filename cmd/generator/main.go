// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/avodah-inc/crossplane-provider-incident-io/config"
	"github.com/crossplane/upjet/pkg/pipeline"
)

func main() {
	// Determine the root directory for code generation output.
	// Default to the current working directory if not specified.
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get working directory: %v\n", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get absolute path: %v\n", err)
		os.Exit(1)
	}

	// Load the provider configuration which includes the Terraform provider
	// schema, resource inclusion list, external name mappings, and all
	// per-domain resource configurations.
	pc := config.GetProvider()

	// Run the Upjet code generation pipeline. This produces:
	// - Go types under apis/ (one package per domain group and version)
	// - CRD YAML under package/crds/
	// - Controllers under internal/controller/
	pipeline.Run(pc, rootDir)
}

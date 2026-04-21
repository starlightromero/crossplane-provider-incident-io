//go:build tools
// +build tools

package tools

// This file imports packages that are used by build tooling but not directly
// by the provider code. It ensures go mod tidy keeps these dependencies.

import (
	_ "github.com/crossplane/crossplane-runtime/apis/common/v1"
	_ "github.com/crossplane/crossplane-tools/cmd/angryjet"
	_ "github.com/crossplane/upjet/pkg/pipeline"
	_ "sigs.k8s.io/controller-runtime/pkg/manager"
)

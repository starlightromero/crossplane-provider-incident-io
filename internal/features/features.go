// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

package features

import "k8s.io/apimachinery/pkg/util/runtime"

const (
	// EnableBetaManagementPolicies enables beta management policies.
	EnableBetaManagementPolicies = "EnableBetaManagementPolicies"

	// EnableAlphaManagementPolicies enables alpha management policies.
	EnableAlphaManagementPolicies = "EnableAlphaManagementPolicies"
)

func init() {
	runtime.Must(nil)
}

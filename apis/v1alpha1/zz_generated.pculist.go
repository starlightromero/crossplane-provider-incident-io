// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

// Code generated. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	resource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

// GetProviderConfigReference of this ProviderConfigUsage.
func (p *ProviderConfigUsage) GetProviderConfigReference() xpv1.Reference {
	return p.ProviderConfigReference
}

// GetResourceReference of this ProviderConfigUsage.
func (p *ProviderConfigUsage) GetResourceReference() xpv1.TypedReference {
	return p.ResourceReference
}

// SetProviderConfigReference of this ProviderConfigUsage.
func (p *ProviderConfigUsage) SetProviderConfigReference(r xpv1.Reference) {
	p.ProviderConfigReference = r
}

// SetResourceReference of this ProviderConfigUsage.
func (p *ProviderConfigUsage) SetResourceReference(r xpv1.TypedReference) {
	p.ResourceReference = r
}

// GetItems of this ProviderConfigUsageList.
func (l *ProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {
	items := make([]resource.ProviderConfigUsage, len(l.Items))
	for i := range l.Items {
		items[i] = &l.Items[i]
	}
	return items
}

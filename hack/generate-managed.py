#!/usr/bin/env python3
"""Generate managed resource method sets for Crossplane types.

This generates the zz_generated.managed.go and zz_generated.managedlist.go
files that angryjet would normally produce. These files delegate the Managed
interface methods from the named Spec/Status fields to the embedded types.
"""

import os

HEADER = """// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

// Code generated. DO NOT EDIT.

package v1alpha1
"""

MANAGED_IMPORTS = """
import (
\txpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
\txpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
)
"""

MANAGEDLIST_IMPORTS = """
import resource "github.com/crossplane/crossplane-runtime/pkg/resource"
"""

def managed_methods(type_name):
    return f"""
// GetCondition of this {type_name}.
func (mg *{type_name}) GetCondition(ct xpv1.ConditionType) xpv1.Condition {{
\treturn mg.Status.GetCondition(ct)
}}

// GetDeletionPolicy of this {type_name}.
func (mg *{type_name}) GetDeletionPolicy() xpv1.DeletionPolicy {{
\treturn mg.Spec.DeletionPolicy
}}

// GetManagementPolicies of this {type_name}.
func (mg *{type_name}) GetManagementPolicies() xpv1.ManagementPolicies {{
\treturn mg.Spec.ManagementPolicies
}}

// GetProviderConfigReference of this {type_name}.
func (mg *{type_name}) GetProviderConfigReference() *xpv1.Reference {{
\treturn mg.Spec.ProviderConfigReference
}}

// GetPublishConnectionDetailsTo of this {type_name}.
func (mg *{type_name}) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {{
\treturn mg.Spec.PublishConnectionDetailsTo
}}

// GetWriteConnectionSecretToReference of this {type_name}.
func (mg *{type_name}) GetWriteConnectionSecretToReference() *xpv1.SecretReference {{
\treturn mg.Spec.WriteConnectionSecretToReference
}}

// SetConditions of this {type_name}.
func (mg *{type_name}) SetConditions(c ...xpv1.Condition) {{
\tmg.Status.SetConditions(c...)
}}

// SetDeletionPolicy of this {type_name}.
func (mg *{type_name}) SetDeletionPolicy(r xpv1.DeletionPolicy) {{
\tmg.Spec.DeletionPolicy = r
}}

// SetManagementPolicies of this {type_name}.
func (mg *{type_name}) SetManagementPolicies(r xpv1.ManagementPolicies) {{
\tmg.Spec.ManagementPolicies = r
}}

// SetProviderConfigReference of this {type_name}.
func (mg *{type_name}) SetProviderConfigReference(r *xpv1.Reference) {{
\tmg.Spec.ProviderConfigReference = r
}}

// SetPublishConnectionDetailsTo of this {type_name}.
func (mg *{type_name}) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {{
\tmg.Spec.PublishConnectionDetailsTo = r
}}

// SetWriteConnectionSecretToReference of this {type_name}.
func (mg *{type_name}) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {{
\tmg.Spec.WriteConnectionSecretToReference = r
}}

// SetObservedGeneration of this {type_name}.
func (mg *{type_name}) SetObservedGeneration(generation int64) {{
\tmg.Status.SetObservedGeneration(generation)
}}

// GetObservedGeneration of this {type_name}.
func (mg *{type_name}) GetObservedGeneration() int64 {{
\treturn mg.Status.GetObservedGeneration()
}}

// Ensure {type_name} implements resource.Managed.
var _ xpresource.Managed = &{type_name}{{}}
"""


def managedlist_methods(type_name, list_name):
    return f"""
// GetItems of this {list_name}.
func (l *{list_name}) GetItems() []resource.Managed {{
\titems := make([]resource.Managed, len(l.Items))
\tfor i := range l.Items {{
\t\titems[i] = &l.Items[i]
\t}}
\treturn items
}}
"""


# Map of domain -> list of (type_name, list_name) pairs
domains = {
    'alerting': [
        ('Attribute', 'AttributeList'),
        ('Route', 'RouteList'),
        ('Source', 'SourceList'),
    ],
    'automation': [
        ('Workflow', 'WorkflowList'),
    ],
    'catalog': [
        ('CatalogType', 'CatalogTypeList'),
        ('CatalogEntries', 'CatalogEntriesList'),
        ('Entry', 'EntryList'),
        ('TypeAttribute', 'TypeAttributeList'),
    ],
    'incident': [
        ('Field', 'FieldList'),
        ('FieldOption', 'FieldOptionList'),
        ('Role', 'RoleList'),
        ('Severity', 'SeverityList'),
        ('Status', 'StatusList'),
    ],
    'oncall': [
        ('Path', 'PathList'),
        ('Schedule', 'ScheduleList'),
        ('Window', 'WindowList'),
    ],
}

for domain, types in domains.items():
    dir_path = f'apis/{domain}/v1alpha1'

    # Generate zz_generated.managed.go
    managed_content = HEADER + MANAGED_IMPORTS
    for type_name, _ in types:
        managed_content += managed_methods(type_name)

    managed_path = os.path.join(dir_path, 'zz_generated.managed.go')
    with open(managed_path, 'w') as f:
        f.write(managed_content)
    print(f'Written {managed_path}')

    # Generate zz_generated.managedlist.go
    list_content = HEADER + MANAGEDLIST_IMPORTS
    for type_name, list_name in types:
        list_content += managedlist_methods(type_name, list_name)

    list_path = os.path.join(dir_path, 'zz_generated.managedlist.go')
    with open(list_path, 'w') as f:
        f.write(list_content)
    print(f'Written {list_path}')

// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

// Code generated. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

// GetCondition of this Path.
func (mg *Path) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Path.
func (mg *Path) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this Path.
func (mg *Path) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this Path.
func (mg *Path) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetWriteConnectionSecretToReference of this Path.
func (mg *Path) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Path.
func (mg *Path) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Path.
func (mg *Path) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this Path.
func (mg *Path) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this Path.
func (mg *Path) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetWriteConnectionSecretToReference of this Path.
func (mg *Path) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this Path.
func (mg *Path) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this Path.
func (mg *Path) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure Path implements resource.Managed.
var _ xpresource.Managed = &Path{}

// GetCondition of this Schedule.
func (mg *Schedule) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Schedule.
func (mg *Schedule) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this Schedule.
func (mg *Schedule) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this Schedule.
func (mg *Schedule) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetWriteConnectionSecretToReference of this Schedule.
func (mg *Schedule) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Schedule.
func (mg *Schedule) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Schedule.
func (mg *Schedule) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this Schedule.
func (mg *Schedule) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this Schedule.
func (mg *Schedule) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetWriteConnectionSecretToReference of this Schedule.
func (mg *Schedule) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this Schedule.
func (mg *Schedule) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this Schedule.
func (mg *Schedule) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure Schedule implements resource.Managed.
var _ xpresource.Managed = &Schedule{}

// GetCondition of this Window.
func (mg *Window) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Window.
func (mg *Window) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this Window.
func (mg *Window) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this Window.
func (mg *Window) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetWriteConnectionSecretToReference of this Window.
func (mg *Window) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Window.
func (mg *Window) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Window.
func (mg *Window) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this Window.
func (mg *Window) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this Window.
func (mg *Window) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetWriteConnectionSecretToReference of this Window.
func (mg *Window) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this Window.
func (mg *Window) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this Window.
func (mg *Window) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure Window implements resource.Managed.
var _ xpresource.Managed = &Window{}

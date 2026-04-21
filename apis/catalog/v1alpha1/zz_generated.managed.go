// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

// Code generated. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
)

// GetCondition of this CatalogType.
func (mg *CatalogType) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this CatalogType.
func (mg *CatalogType) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this CatalogType.
func (mg *CatalogType) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this CatalogType.
func (mg *CatalogType) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this CatalogType.
func (mg *CatalogType) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this CatalogType.
func (mg *CatalogType) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this CatalogType.
func (mg *CatalogType) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this CatalogType.
func (mg *CatalogType) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this CatalogType.
func (mg *CatalogType) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this CatalogType.
func (mg *CatalogType) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this CatalogType.
func (mg *CatalogType) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this CatalogType.
func (mg *CatalogType) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this CatalogType.
func (mg *CatalogType) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this CatalogType.
func (mg *CatalogType) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure CatalogType implements resource.Managed.
var _ xpresource.Managed = &CatalogType{}

// GetCondition of this CatalogEntries.
func (mg *CatalogEntries) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this CatalogEntries.
func (mg *CatalogEntries) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this CatalogEntries.
func (mg *CatalogEntries) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this CatalogEntries.
func (mg *CatalogEntries) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this CatalogEntries.
func (mg *CatalogEntries) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this CatalogEntries.
func (mg *CatalogEntries) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this CatalogEntries.
func (mg *CatalogEntries) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this CatalogEntries.
func (mg *CatalogEntries) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this CatalogEntries.
func (mg *CatalogEntries) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this CatalogEntries.
func (mg *CatalogEntries) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this CatalogEntries.
func (mg *CatalogEntries) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this CatalogEntries.
func (mg *CatalogEntries) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this CatalogEntries.
func (mg *CatalogEntries) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this CatalogEntries.
func (mg *CatalogEntries) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure CatalogEntries implements resource.Managed.
var _ xpresource.Managed = &CatalogEntries{}

// GetCondition of this Entry.
func (mg *Entry) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Entry.
func (mg *Entry) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this Entry.
func (mg *Entry) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this Entry.
func (mg *Entry) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this Entry.
func (mg *Entry) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this Entry.
func (mg *Entry) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Entry.
func (mg *Entry) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Entry.
func (mg *Entry) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this Entry.
func (mg *Entry) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this Entry.
func (mg *Entry) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this Entry.
func (mg *Entry) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this Entry.
func (mg *Entry) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this Entry.
func (mg *Entry) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this Entry.
func (mg *Entry) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure Entry implements resource.Managed.
var _ xpresource.Managed = &Entry{}

// GetCondition of this TypeAttribute.
func (mg *TypeAttribute) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this TypeAttribute.
func (mg *TypeAttribute) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this TypeAttribute.
func (mg *TypeAttribute) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this TypeAttribute.
func (mg *TypeAttribute) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this TypeAttribute.
func (mg *TypeAttribute) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this TypeAttribute.
func (mg *TypeAttribute) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this TypeAttribute.
func (mg *TypeAttribute) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this TypeAttribute.
func (mg *TypeAttribute) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this TypeAttribute.
func (mg *TypeAttribute) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this TypeAttribute.
func (mg *TypeAttribute) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this TypeAttribute.
func (mg *TypeAttribute) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this TypeAttribute.
func (mg *TypeAttribute) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// SetObservedGeneration of this TypeAttribute.
func (mg *TypeAttribute) SetObservedGeneration(generation int64) {
	mg.Status.SetObservedGeneration(generation)
}

// GetObservedGeneration of this TypeAttribute.
func (mg *TypeAttribute) GetObservedGeneration() int64 {
	return mg.Status.GetObservedGeneration()
}

// Ensure TypeAttribute implements resource.Managed.
var _ xpresource.Managed = &TypeAttribute{}

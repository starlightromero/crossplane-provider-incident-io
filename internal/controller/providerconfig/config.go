// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

package providerconfig

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/providerconfig"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/controller"

	"github.com/avodah-inc/crossplane-provider-incident-io/apis/v1alpha1"
)

// Setup adds a controller that reconciles ProviderConfigs.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := providerconfig.ControllerName(v1alpha1.ProviderConfigGroupVersionKind.Kind)

	of := resource.ProviderConfigKinds{
		Config:    v1alpha1.ProviderConfigGroupVersionKind,
		Usage:     v1alpha1.ProviderConfigUsageGroupVersionKind,
		UsageList: v1alpha1.ProviderConfigUsageListGroupVersionKind,
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.ProviderConfig{}).
		Watches(&v1alpha1.ProviderConfigUsage{}, &resource.EnqueueRequestForProviderConfig{}).
		Complete(providerconfig.NewReconciler(mgr, of, providerconfig.WithLogger(o.Logger.WithValues("controller", name)),
			providerconfig.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

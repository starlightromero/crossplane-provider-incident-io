// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	attribute "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/alerting/attribute"
	route "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/alerting/route"
	source "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/alerting/source"
	workflow "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/automation/workflow"
	catalogentries "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/catalog/catalogentries"
	catalogtype "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/catalog/catalogtype"
	entry "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/catalog/entry"
	typeattribute "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/catalog/typeattribute"
	field "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/incident/field"
	fieldoption "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/incident/fieldoption"
	role "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/incident/role"
	severity "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/incident/severity"
	status "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/incident/status"
	path "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/oncall/path"
	schedule "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/oncall/schedule"
	window "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/oncall/window"
	providerconfig "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		attribute.Setup,
		route.Setup,
		source.Setup,
		workflow.Setup,
		catalogentries.Setup,
		catalogtype.Setup,
		entry.Setup,
		typeattribute.Setup,
		field.Setup,
		fieldoption.Setup,
		role.Setup,
		severity.Setup,
		status.Setup,
		path.Setup,
		schedule.Setup,
		window.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

package config

import ujconfig "github.com/crossplane/upjet/pkg/config"

// ExternalNameConfigs contains all external name configurations for this
// provider. All Incident.io resources use server-assigned UUIDs, so
// IdentifierFromProvider is the appropriate strategy across the board.
var ExternalNameConfigs = map[string]ujconfig.ExternalName{
	// Alerting domain
	"incident_alert_attribute": ujconfig.IdentifierFromProvider,
	"incident_alert_route":     ujconfig.IdentifierFromProvider,
	"incident_alert_source":    ujconfig.IdentifierFromProvider,

	// Catalog domain
	"incident_catalog_type":           ujconfig.IdentifierFromProvider,
	"incident_catalog_type_attribute": ujconfig.IdentifierFromProvider,
	"incident_catalog_entry":          ujconfig.IdentifierFromProvider,
	"incident_catalog_entries":        ujconfig.IdentifierFromProvider,

	// Incident configuration domain
	"incident_custom_field":        ujconfig.IdentifierFromProvider,
	"incident_custom_field_option": ujconfig.IdentifierFromProvider,
	"incident_incident_role":       ujconfig.IdentifierFromProvider,
	"incident_severity":            ujconfig.IdentifierFromProvider,
	"incident_status":              ujconfig.IdentifierFromProvider,

	// On-call domain
	"incident_escalation_path":    ujconfig.IdentifierFromProvider,
	"incident_schedule":           ujconfig.IdentifierFromProvider,
	"incident_maintenance_window": ujconfig.IdentifierFromProvider,

	// Automation domain
	"incident_workflow": ujconfig.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the external name of those resources.
func ExternalNameConfigurations() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually. This is used as the include list for the provider.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}

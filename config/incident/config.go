package incident

import ujconfig "github.com/crossplane/upjet/pkg/config"

const group = "incident"

// Configure adds configurations for the incident configuration domain resources.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("incident_custom_field", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_custom_field_option", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_incident_role", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_severity", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_status", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
}

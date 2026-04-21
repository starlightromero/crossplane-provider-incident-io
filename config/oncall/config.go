package oncall

import ujconfig "github.com/crossplane/upjet/pkg/config"

const group = "oncall"

// Configure adds configurations for the on-call domain resources.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("incident_escalation_path", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_schedule", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_maintenance_window", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
}

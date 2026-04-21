package alerting

import ujconfig "github.com/crossplane/upjet/pkg/config"

const group = "alerting"

// Configure adds configurations for the alerting domain resources.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("incident_alert_attribute", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_alert_route", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_alert_source", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
}

package automation

import ujconfig "github.com/crossplane/upjet/pkg/config"

const group = "automation"

// Configure adds configurations for the automation domain resources.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("incident_workflow", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
}

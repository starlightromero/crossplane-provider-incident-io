package catalog

import ujconfig "github.com/crossplane/upjet/v2/pkg/config"

const group = "catalog"

// Configure adds configurations for the catalog domain resources.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("incident_catalog_type", func(r *ujconfig.Resource) {
		r.ShortGroup = group
		// Override Kind to avoid "type" which is a Go reserved keyword
		// in the generated controller package name.
		r.Kind = "CatalogType"
	})
	p.AddResourceConfigurator("incident_catalog_type_attribute", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_catalog_entry", func(r *ujconfig.Resource) {
		r.ShortGroup = group
	})
	p.AddResourceConfigurator("incident_catalog_entries", func(r *ujconfig.Resource) {
		r.ShortGroup = group
		// Override Kind to avoid plural collision with CatalogEntry
		// (both "Entry" and "Entries" pluralize to "entries").
		r.Kind = "CatalogEntries"
	})
}

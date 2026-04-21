package config

import (
	// Note(provider): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/avodah-inc/crossplane-provider-incident-io/config/alerting"
	"github.com/avodah-inc/crossplane-provider-incident-io/config/automation"
	"github.com/avodah-inc/crossplane-provider-incident-io/config/catalog"
	"github.com/avodah-inc/crossplane-provider-incident-io/config/incident"
	"github.com/avodah-inc/crossplane-provider-incident-io/config/oncall"
)

const (
	// ResourcePrefix is the Terraform resource prefix for the incident-io provider.
	ResourcePrefix = "incident"

	// ModulePath is the Go module path for this provider.
	ModulePath = "github.com/avodah-inc/crossplane-provider-incident-io"

	// RootGroup is the root API group for all managed resources.
	RootGroup = "incidentio.crossplane.io"

	// ShortName is the short name for the provider.
	ShortName = "incident-io"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns the provider configuration for the incident-io provider.
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), ResourcePrefix, ModulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup(RootGroup),
		ujconfig.WithShortName(ShortName),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
	)

	for _, configure := range []func(provider *ujconfig.Provider){
		alerting.Configure,
		catalog.Configure,
		incident.Configure,
		oncall.Configure,
		automation.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()

	return pc
}

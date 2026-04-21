package clients

import (
	"context"
	"strings"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/upjet/pkg/terraform"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/avodah-inc/crossplane-provider-incident-io/apis/v1alpha1"
)

const (
	// DefaultEndpoint is the default Incident.io API endpoint.
	DefaultEndpoint = "https://api.incident.io"

	// Terraform provider configuration keys.
	keyAPIKey   = "api_key"
	keyEndpoint = "endpoint"

	errNoProviderConfig    = "cannot get referenced ProviderConfig"
	errGetCredentials      = "cannot get credentials"
	errExtractCredentials  = "cannot extract credentials from secret"
	errCredentialsEmpty    = "extracted credentials are empty"
	errTrackProviderConfig = "cannot track ProviderConfig usage"
)

// TerraformSetupBuilder returns a terraform.SetupFn that reads credentials from
// the ProviderConfig's referenced Secret and builds the Terraform provider
// configuration for the incident-io provider.
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		setup := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return setup, errors.New(errNoProviderConfig)
		}

		pc := &v1alpha1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return setup, errors.Wrap(err, errNoProviderConfig)
		}

		// Extract credentials from the referenced Secret.
		creds, err := resource.CommonCredentialExtractor(
			ctx,
			pc.Spec.Credentials.Source,
			client,
			pc.Spec.Credentials.CommonCredentialSelectors,
		)
		if err != nil {
			return setup, errors.Wrap(err, errExtractCredentials)
		}

		apiKey := strings.TrimSpace(string(creds))
		if apiKey == "" {
			return setup, errors.New(errCredentialsEmpty)
		}

		// Build provider configuration map.
		setup.Configuration = terraform.ProviderConfiguration{
			keyAPIKey:   apiKey,
			keyEndpoint: DefaultEndpoint,
		}

		return setup, nil
	}
}

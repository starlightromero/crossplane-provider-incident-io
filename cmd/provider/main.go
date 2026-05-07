// SPDX-FileCopyrightText: 2024 Avodah Inc.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/avodah-inc/crossplane-provider-incident-io/apis"
	"github.com/avodah-inc/crossplane-provider-incident-io/config"
	"github.com/avodah-inc/crossplane-provider-incident-io/internal/clients"
	intcontroller "github.com/avodah-inc/crossplane-provider-incident-io/internal/controller"
)

// Version is set via ldflags at build time.
var Version = "dev"

const (
	providerSource  = "incident-io/incident"
	providerVersion = ">=0.0.1"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]),
			"Crossplane provider for Incident.io").DefaultEnvars()

		debug = app.Flag("debug",
			"Run with debug logging.").Short('d').Bool()

		pollInterval = app.Flag("poll",
			"How often individual resources will be checked for drift from the desired state.").
			Default("10m").Duration()

		maxReconcileRate = app.Flag("max-reconcile-rate",
			"The global maximum rate per second at which resources may be checked for drift from the desired state.").
			Default("10").Int()

		leaderElection = app.Flag("leader-election",
			"Use leader election for the controller manager.").
			Short('l').Envar("LEADER_ELECTION").Default("false").Bool()

		enableManagementPolicies = app.Flag("enable-management-policies",
			"Enable support for Management Policies.").
			Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
	)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-incident-io"))

	if *debug {
		ctrl.SetLogger(zl)
	}

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:             *leaderElection,
		LeaderElectionID:           "crossplane-leader-election-provider-incident-io",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add APIs to scheme")

	featureFlags := &feature.Flags{}
	if *enableManagementPolicies {
		featureFlags.Enable(feature.EnableBetaManagementPolicies)
	}

	// Set up the workspace store for Terraform state management.
	ws := terraform.NewWorkspaceStore(log)

	// Set up the operation tracker store for async operations.
	tracker := tjcontroller.NewOperationStore(log)

	o := tjcontroller.Options{
		Provider: config.GetProvider(),
		Options: controller.Options{
			Logger:                  log,
			MaxConcurrentReconciles: *maxReconcileRate,
			PollInterval:            *pollInterval,
			GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
			Features:                featureFlags,
		},
		WorkspaceStore:        ws,
		SetupFn:               clients.TerraformSetupBuilder(Version, providerSource, providerVersion),
		PollJitter:            time.Second,
		OperationTrackerStore: tracker,
	}

	kingpin.FatalIfError(intcontroller.Setup(mgr, o), "Cannot setup controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

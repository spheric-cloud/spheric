// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	computescheduler "spheric.cloud/spheric/internal/controllers/core/scheduler"

	"k8s.io/utils/lru"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	coreclient "spheric.cloud/spheric/internal/client/core"
	corecontrollers "spheric.cloud/spheric/internal/controllers/core"
	certificatespheric "spheric.cloud/spheric/internal/controllers/core/certificate/spheric"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/ironcore-dev/controller-utils/cmdutils/switches"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	instanceEphemeralVolumeController = "machineephemeralvolume"
	instanceSchedulerController       = "machinescheduler"
	instanceTypeController            = "instancetype"
	diskReleaseController             = "volumerelease"
	networkProtectionController       = "networkprotection"
	certificateApprovalController     = "certificateapproval"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var prefixAllocationTimeout time.Duration
	var volumeBindTimeout time.Duration
	var virtualIPBindTimeout time.Duration
	var networkInterfaceBindTimeout time.Duration
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.DurationVar(&prefixAllocationTimeout, "prefix-allocation-timeout", 1*time.Second, "Time to wait until considering a pending allocation failed.")
	flag.DurationVar(&volumeBindTimeout, "disk-bind-timeout", 10*time.Second, "Time to wait until considering a disk bind to be failed.")
	flag.DurationVar(&virtualIPBindTimeout, "virtual-ip-bind-timeout", 10*time.Second, "Time to wait until considering a virtual ip bind to be failed.")
	flag.DurationVar(&networkInterfaceBindTimeout, "network-interface-bind-timeout", 10*time.Second, "Time to wait until considering a network interface bind to be failed.")

	controllers := switches.New(
		instanceEphemeralVolumeController,
		instanceSchedulerController,
		instanceTypeController,
		diskReleaseController,
		networkProtectionController,
		certificateApprovalController,
	)
	flag.Var(controllers, "controllers",
		fmt.Sprintf("Controllers to enable. All controllers: %v. Disabled-by-default controllers: %v",
			controllers.All(),
			controllers.DisabledByDefault(),
		),
	)

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)
	ctx := ctrl.SetupSignalHandler()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Logger:                 logger,
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "d0ae00be.spheric.cloud",
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Register controllers

	if controllers.Enabled(instanceEphemeralVolumeController) {
		if err := (&corecontrollers.InstanceEphemeralDiskReconciler{
			Client: mgr.GetClient(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "InstanceEphemeralVolume")
			os.Exit(1)
		}
	}

	if controllers.Enabled(instanceSchedulerController) {
		schedulerCache := computescheduler.NewCache(mgr.GetLogger(), computescheduler.DefaultCacheStrategy)
		if err := mgr.Add(schedulerCache); err != nil {
			setupLog.Error(err, "unable to create cache", "controller", "InstanceSchedulerCache")
			os.Exit(1)
		}

		if err := (&corecontrollers.InstanceScheduler{
			Client:        mgr.GetClient(),
			EventRecorder: mgr.GetEventRecorderFor("instance-scheduler"),
			Cache:         schedulerCache,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "InstanceScheduler")
			os.Exit(1)
		}
	}

	if controllers.Enabled(instanceTypeController) {
		if err := (&corecontrollers.InstanceTypeReconciler{
			Client:    mgr.GetClient(),
			APIReader: mgr.GetAPIReader(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "InstanceType")
			os.Exit(1)
		}
	}

	if controllers.Enabled(diskReleaseController) {
		if err := (&corecontrollers.DiskReleaseReconciler{
			Client:       mgr.GetClient(),
			APIReader:    mgr.GetAPIReader(),
			AbsenceCache: lru.New(500),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DiskRelease")
			os.Exit(1)
		}
	}

	if controllers.Enabled(networkProtectionController) {
		if err := (&corecontrollers.NetworkProtectionReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "NetworkProtection")
			os.Exit(1)
		}
	}

	if controllers.Enabled(certificateApprovalController) {
		if err := (&corecontrollers.CertificateApprovalReconciler{
			Client:      mgr.GetClient(),
			Recognizers: certificatespheric.Recognizers,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CertificateApproval")
			os.Exit(1)
		}
	}

	if controllers.AnyEnabled(instanceEphemeralVolumeController) {
		if err := coreclient.SetupInstanceSpecDiskNamesFieldIndexer(ctx, mgr.GetFieldIndexer()); err != nil {
			setupLog.Error(err, "unable to index field", "field", coreclient.InstanceSpecDiskNamesField)
			os.Exit(1)
		}
	}

	if controllers.AnyEnabled(instanceSchedulerController) {
		if err := coreclient.SetupInstanceSpecFleetRefNameFieldIndexer(ctx, mgr.GetFieldIndexer()); err != nil {
			setupLog.Error(err, "unable to index field", "field", coreclient.InstanceSpecFleetRefNameField)
			os.Exit(1)
		}
	}

	if controllers.AnyEnabled(instanceTypeController) {
		if err := coreclient.SetupInstanceSpecInstanceTypeRefNameFieldIndexer(ctx, mgr.GetFieldIndexer()); err != nil {
			setupLog.Error(err, "unable to setup field indexer", "field", coreclient.InstanceSpecInstanceTypeRefNameField)
			os.Exit(1)
		}
	}

	// healthz / readyz setup

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	metricserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	coreclient "spheric.cloud/spheric/internal/client/core"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	sphereletclient "spheric.cloud/spheric/spherelet/client"
	"spheric.cloud/spheric/spherelet/controllers"
	"spheric.cloud/spheric/spherelet/event/instanceevent"
	"spheric.cloud/spheric/spherelet/event/runtimeevent"
	. "spheric.cloud/spheric/utils/testing"
	"spheric.cloud/spheric/utils/testing/record"

	"spheric.cloud/spheric/spherelet/iri/remote/fake"

	"github.com/ironcore-dev/controller-utils/buildutils"
	"github.com/ironcore-dev/controller-utils/modutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	utilsenvtest "spheric.cloud/spheric/utils/envtest"
	"spheric.cloud/spheric/utils/envtest/apiserver"
	"spheric.cloud/spheric/utils/envtest/controllermanager"
	"spheric.cloud/spheric/utils/envtest/process"
)

var (
	cfg        *rest.Config
	testEnv    *envtest.Environment
	testEnvExt *utilsenvtest.EnvironmentExtensions
	k8sClient  = NewClientPromise()
	srv        *fake.FakeRuntimeService
)

const (
	eventuallyTimeout    = 3 * time.Second
	pollingInterval      = 50 * time.Millisecond
	consistentlyDuration = 1 * time.Second
	apiServiceTimeout    = 1 * time.Minute

	controllerManagerService = "controller-manager"

	fleetName = "test-fleet"

	fooDownwardAPILabel = "custom-downward-api-label"
	fooAnnotation       = "foo"
)

func TestControllers(t *testing.T) {
	SetDefaultConsistentlyPollingInterval(pollingInterval)
	SetDefaultEventuallyPollingInterval(pollingInterval)
	SetDefaultEventuallyTimeout(eventuallyTimeout)
	SetDefaultConsistentlyDuration(consistentlyDuration)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	var err error
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}
	testEnvExt = &utilsenvtest.EnvironmentExtensions{
		APIServiceDirectoryPaths: []string{
			modutils.Dir("spheric.cloud/spheric", "config", "apiserver", "apiservice", "bases"),
		},
		ErrorIfAPIServicePathIsMissing: true,
		AdditionalServices: []utilsenvtest.AdditionalService{
			{
				Name: controllerManagerService,
			},
		},
	}

	cfg, err = utilsenvtest.StartWithExtensions(testEnv, testEnvExt)
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	DeferCleanup(utilsenvtest.StopWithExtensions, testEnv, testEnvExt)

	Expect(corev1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())

	// Init package-level k8sClient
	Expect(k8sClient.FulfillWith(client.New(cfg, client.Options{Scheme: scheme.Scheme}))).To(Succeed())
	SetClient(k8sClient)

	apiSrv, err := apiserver.New(cfg, apiserver.Options{
		MainPath:     "spheric.cloud/spheric/cmd/apiserver",
		BuildOptions: []buildutils.BuildOption{buildutils.ModModeMod},
		ETCDServers:  []string{testEnv.ControlPlane.Etcd.URL.String()},
		Host:         testEnvExt.APIServiceInstallOptions.LocalServingHost,
		Port:         testEnvExt.APIServiceInstallOptions.LocalServingPort,
		CertDir:      testEnvExt.APIServiceInstallOptions.LocalServingCertDir,
	})
	Expect(err).NotTo(HaveOccurred())

	Expect(apiSrv.Start()).To(Succeed())
	DeferCleanup(apiSrv.Stop)

	Expect(utilsenvtest.WaitUntilAPIServicesReadyWithTimeout(apiServiceTimeout, testEnvExt, k8sClient, scheme.Scheme)).To(Succeed())

	ctrlMgr, err := controllermanager.New(cfg, controllermanager.Options{
		Args:         process.EmptyArgs().Set("controllers", "*"),
		MainPath:     "spheric.cloud/spheric/cmd/controller-manager",
		BuildOptions: []buildutils.BuildOption{buildutils.ModModeMod},
		Host:         testEnvExt.GetAdditionalServiceHost(controllerManagerService),
		Port:         testEnvExt.GetAdditionalServicePort(controllerManagerService),
	})
	Expect(err).NotTo(HaveOccurred())

	Expect(ctrlMgr.Start()).To(Succeed())
	DeferCleanup(ctrlMgr.Stop)

	mgrCtx, cancel := context.WithCancel(context.Background())
	DeferCleanup(cancel)

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		Metrics: metricserver.Options{
			BindAddress: "0",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	indexer := k8sManager.GetFieldIndexer()
	Expect(sphereletclient.SetupInstanceSpecDiskNamesField(mgrCtx, indexer, fleetName)).To(Succeed())
	Expect(sphereletclient.SetupInstanceSpecSecretNamesField(mgrCtx, indexer, fleetName)).To(Succeed())
	Expect(coreclient.SetupInstanceSpecFleetRefNameFieldIndexer(mgrCtx, indexer)).To(Succeed())

	onInitialized := func() {
		srv = fake.NewFakeRuntimeService()

		Expect((&controllers.InstanceReconciler{
			EventRecorder:          &record.LogRecorder{Logger: GinkgoLogr},
			Client:                 k8sManager.GetClient(),
			InstanceRuntime:        srv,
			InstanceRuntimeName:    fake.RuntimeName,
			InstanceRuntimeVersion: fake.Version,
			FleetName:              fleetName,
			DownwardAPILabels: map[string]string{
				fooDownwardAPILabel: fmt.Sprintf("metadata.annotations['%s']", fooAnnotation),
			},
		}).SetupWithManager(k8sManager)).To(Succeed())

		instanceEvents := instanceevent.NewGenerator(func(ctx context.Context) ([]*iri.Instance, error) {
			res, err := srv.ListInstances(ctx, &iri.ListInstancesRequest{})
			if err != nil {
				return nil, err
			}
			return res.Instances, nil
		}, instanceevent.GeneratorOptions{})

		Expect(k8sManager.Add(instanceEvents)).To(Succeed())

		runtimeEvents := runtimeevent.NewGenerator(func(ctx context.Context) (*iri.RuntimeResources, error) {
			res, err := srv.Status(ctx, &iri.StatusRequest{})
			if err != nil {
				return nil, err
			}
			return res.Allocatable, nil
		}, runtimeevent.GeneratorOptions{})

		Expect(k8sManager.Add(runtimeEvents)).To(Succeed())

		Expect((&controllers.InstanceAnnotatorReconciler{
			Client:         k8sManager.GetClient(),
			InstanceEvents: instanceEvents,
		}).SetupWithManager(k8sManager)).To(Succeed())

		Expect((&controllers.FleetReconciler{
			Client:          k8sManager.GetClient(),
			InstanceRuntime: srv,
			FleetName:       fleetName,
		}).SetupWithManager(k8sManager)).To(Succeed())

		Expect((&controllers.FleetAnnotatorReconciler{
			Client:        k8sManager.GetClient(),
			FleetName:     fleetName,
			RuntimeEvents: runtimeEvents,
		}).SetupWithManager(k8sManager)).To(Succeed())
	}

	initialized := make(chan struct{})
	fleetInit := controllers.FleetInit{
		Client:     k8sManager.GetClient(),
		FleetName:  fleetName,
		ProviderID: "my-provider-id",
		OnInitialized: func(ctx context.Context) error {
			defer GinkgoRecover()

			onInitialized()
			close(initialized)
			return nil
		},
		OnFailed: func(ctx context.Context, reason error) error {
			defer GinkgoRecover()
			Fail(reason.Error())
			return nil
		},
	}
	Expect(fleetInit.SetupWithManager(k8sManager)).To(Succeed())

	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(mgrCtx)).To(Succeed(), "failed to start manager")
	}()

	Eventually(initialized).Should(BeClosed(), "did not successfully initialize")
})

func SetupInstanceType() *corev1alpha1.InstanceType {
	return SetupNewObject(k8sClient, func(typ *corev1alpha1.InstanceType) {
		*typ = corev1alpha1.InstanceType{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-type-",
			},
			Capabilities: corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:    resource.MustParse("1"),
				corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
	})
}

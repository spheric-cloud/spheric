// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/lru"
	coreclient "spheric.cloud/spheric/internal/client/core"
	"spheric.cloud/spheric/internal/controllers/core/scheduler"
	. "spheric.cloud/spheric/utils/testing"

	"github.com/ironcore-dev/controller-utils/buildutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/internal/controllers/core"
	certificatespheric "spheric.cloud/spheric/internal/controllers/core/certificate/spheric"
	utilsenvtest "spheric.cloud/spheric/utils/envtest"
	"spheric.cloud/spheric/utils/envtest/apiserver"
)

const (
	pollingInterval      = 50 * time.Millisecond
	eventuallyTimeout    = 3 * time.Second
	consistentlyDuration = 1 * time.Second
	apiServiceTimeout    = 5 * time.Minute
)

var (
	k8sClient  = NewClientPromise()
	testEnv    *envtest.Environment
	testEnvExt *utilsenvtest.EnvironmentExtensions
)

func TestCore(t *testing.T) {
	SetDefaultConsistentlyPollingInterval(pollingInterval)
	SetDefaultEventuallyPollingInterval(pollingInterval)
	SetDefaultEventuallyTimeout(eventuallyTimeout)
	SetDefaultConsistentlyDuration(consistentlyDuration)
	RegisterFailHandler(Fail)

	RunSpecs(t, "Core Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}
	testEnvExt = &utilsenvtest.EnvironmentExtensions{
		APIServiceDirectoryPaths:       []string{filepath.Join("..", "..", "..", "config", "apiserver", "apiservice", "bases")},
		ErrorIfAPIServicePathIsMissing: true,
	}

	cfg, err := utilsenvtest.StartWithExtensions(testEnv, testEnvExt)
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	DeferCleanup(utilsenvtest.StopWithExtensions, testEnv, testEnvExt)

	Expect(corev1alpha1.AddToScheme(scheme.Scheme)).Should(Succeed())

	// Init package-level k8sClient
	Expect(k8sClient.FulfillWith(client.New(cfg, client.Options{Scheme: scheme.Scheme}))).To(Succeed())
	komega.SetClient(k8sClient)

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

	ctx, cancel := context.WithCancel(context.Background())
	DeferCleanup(cancel)

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		Metrics: metricserver.Options{
			BindAddress: "0",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(coreclient.SetupInstanceSpecDiskNamesFieldIndexer(ctx, k8sManager.GetFieldIndexer())).To(Succeed())
	Expect(coreclient.SetupInstanceSpecFleetRefNameFieldIndexer(ctx, k8sManager.GetFieldIndexer())).To(Succeed())
	Expect(coreclient.SetupInstanceSpecInstanceTypeRefNameFieldIndexer(ctx, k8sManager.GetFieldIndexer())).To(Succeed())
	Expect(coreclient.SetupSubnetSpecNetworkRefNameField(ctx, k8sManager.GetFieldIndexer())).To(Succeed())

	schedulerCache := scheduler.NewCache(k8sManager.GetLogger(), scheduler.DefaultCacheStrategy)
	Expect(k8sManager.Add(schedulerCache)).To(Succeed())

	Expect((&core.CertificateApprovalReconciler{
		Client:      k8sManager.GetClient(),
		Recognizers: certificatespheric.Recognizers,
	}).SetupWithManager(k8sManager)).To(Succeed())

	Expect((&core.DiskReleaseReconciler{
		Client:       k8sManager.GetClient(),
		APIReader:    k8sManager.GetAPIReader(),
		AbsenceCache: lru.New(128),
	}).SetupWithManager(k8sManager)).To(Succeed())

	Expect((&core.InstanceEphemeralDiskReconciler{
		Client: k8sManager.GetClient(),
	}).SetupWithManager(k8sManager)).To(Succeed())

	Expect((&core.InstanceTypeReconciler{
		Client:    k8sManager.GetClient(),
		APIReader: k8sManager.GetAPIReader(),
	}).SetupWithManager(k8sManager)).To(Succeed())

	Expect((&core.NetworkProtectionReconciler{
		Client: k8sManager.GetClient(),
		Scheme: scheme.Scheme,
	}).SetupWithManager(k8sManager)).To(Succeed())

	Expect((&core.InstanceScheduler{
		EventRecorder: &record.FakeRecorder{},
		Client:        k8sManager.GetClient(),
		Cache:         schedulerCache,
	}).SetupWithManager(k8sManager)).To(Succeed())

	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(ctx)).To(Succeed(), "failed to start manager")
	}()
})

func SetupInstanceType() *corev1alpha1.InstanceType {
	return SetupNewObject(k8sClient, func(instanceType *corev1alpha1.InstanceType) {
		*instanceType = corev1alpha1.InstanceType{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "instance-type-",
			},
			Capabilities: corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:    resource.MustParse("1"),
				corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
	})
}

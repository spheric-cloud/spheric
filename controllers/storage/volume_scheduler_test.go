// Copyright 2021 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	storagev1alpha1 "github.com/onmetal/onmetal-api/apis/storage/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("VolumeScheduler", func() {
	ctx := ctrl.SetupSignalHandler()
	ns := SetupTest(ctx)

	It("should schedule volumes on storage pools", func() {
		By("creating a storage pool")
		storagePool := &storagev1alpha1.StoragePool{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-pool-",
			},
		}
		Expect(k8sClient.Create(ctx, storagePool)).To(Succeed(), "failed to create storage pool")

		By("patching the storage pool status to contain a storage class")
		storagePoolBase := storagePool.DeepCopy()
		storagePool.Status.AvailableStorageClasses = []corev1.LocalObjectReference{{Name: "my-volumeclass"}}
		Expect(k8sClient.Status().Patch(ctx, storagePool, client.MergeFrom(storagePoolBase))).
			To(Succeed(), "failed to patch storage pool status")

		By("creating a volume w/ the requested storage class")
		volume := &storagev1alpha1.Volume{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-volume-",
			},
			Spec: storagev1alpha1.VolumeSpec{
				StorageClass: corev1.LocalObjectReference{
					Name: "my-volumeclass",
				},
			},
		}
		Expect(k8sClient.Create(ctx, volume)).To(Succeed(), "failed to create volume")

		By("waiting for the volume to be scheduled onto the storage pool")
		volumeKey := client.ObjectKeyFromObject(volume)
		Eventually(func(g Gomega) {
			Expect(k8sClient.Get(ctx, volumeKey, volume)).To(Succeed(), "failed to get volume")
			g.Expect(volume.Spec.StoragePool.Name).To(Equal(storagePool.Name))
			g.Expect(volume.Status.State).To(Equal(storagev1alpha1.VolumeStatePending))
		}).Should(Succeed())
	})

	It("should schedule schedule volumes onto storage pools if the pool becomes available later than the volume", func() {
		By("creating a volume w/ the requested storage class")
		volume := &storagev1alpha1.Volume{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-volume-",
			},
			Spec: storagev1alpha1.VolumeSpec{
				StorageClass: corev1.LocalObjectReference{
					Name: "my-volumeclass",
				},
			},
		}
		Expect(k8sClient.Create(ctx, volume)).To(Succeed(), "failed to create volume")

		By("waiting for the volume to indicate it is pending")
		volumeKey := client.ObjectKeyFromObject(volume)
		Eventually(func(g Gomega) {
			Expect(k8sClient.Get(ctx, volumeKey, volume)).To(Succeed())
			g.Expect(volume.Spec.StoragePool.Name).To(BeEmpty())
			g.Expect(volume.Status.State).To(Equal(storagev1alpha1.VolumeStatePending))
		}).Should(Succeed())

		By("creating a storage pool")
		storagePool := &storagev1alpha1.StoragePool{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-pool-",
			},
		}
		Expect(k8sClient.Create(ctx, storagePool)).To(Succeed(), "failed to create storage pool")

		By("patching the storage pool status to contain a storage class")
		storagePoolBase := storagePool.DeepCopy()
		storagePool.Status.AvailableStorageClasses = []corev1.LocalObjectReference{{Name: "my-volumeclass"}}
		Expect(k8sClient.Status().Patch(ctx, storagePool, client.MergeFrom(storagePoolBase))).
			To(Succeed(), "failed to patch storage pool status")

		By("waiting for the volume to be scheduled onto the storage pool")
		Eventually(func() string {
			Expect(k8sClient.Get(ctx, volumeKey, volume)).To(Succeed(), "failed to get volume")
			return volume.Spec.StoragePool.Name
		}).Should(Equal(storagePool.Name))
	})

	It("should schedule onto storage pools with matching labels", func() {
		By("creating a storage pool w/o matching labels")
		storagePoolNoMatchingLabels := &storagev1alpha1.StoragePool{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-pool-",
			},
		}
		Expect(k8sClient.Create(ctx, storagePoolNoMatchingLabels)).To(Succeed(), "failed to create storage pool")

		By("patching the storage pool status to contain a storage class")
		storagePoolNoMatchingLabelsBase := storagePoolNoMatchingLabels.DeepCopy()
		storagePoolNoMatchingLabels.Status.AvailableStorageClasses = []corev1.LocalObjectReference{{Name: "my-storageclass"}}
		Expect(k8sClient.Status().Patch(ctx, storagePoolNoMatchingLabels, client.MergeFrom(storagePoolNoMatchingLabelsBase))).
			To(Succeed(), "failed to patch storage pool status")

		By("creating a storage pool w/ matching labels")
		storagePoolMatchingLabels := &storagev1alpha1.StoragePool{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-pool-",
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		}
		Expect(k8sClient.Create(ctx, storagePoolMatchingLabels)).To(Succeed(), "failed to create storage pool")

		By("patching the storage pool status to contain a storage class")
		storagePoolMatchingLabelsBase := storagePoolMatchingLabels.DeepCopy()
		storagePoolMatchingLabels.Status.AvailableStorageClasses = []corev1.LocalObjectReference{{Name: "my-storageclass"}}
		Expect(k8sClient.Status().Patch(ctx, storagePoolMatchingLabels, client.MergeFrom(storagePoolMatchingLabelsBase))).
			To(Succeed(), "failed to patch storage pool status")

		By("creating a volume w/ the requested storage class")
		volume := &storagev1alpha1.Volume{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-volume-",
			},
			Spec: storagev1alpha1.VolumeSpec{
				StoragePoolSelector: map[string]string{
					"foo": "bar",
				},
				StorageClass: corev1.LocalObjectReference{
					Name: "my-storageclass",
				},
			},
		}
		Expect(k8sClient.Create(ctx, volume)).To(Succeed(), "failed to create volume")

		By("waiting for the volume to be scheduled onto the storage pool")
		volumeKey := client.ObjectKeyFromObject(volume)
		Eventually(func(g Gomega) {
			Expect(k8sClient.Get(ctx, volumeKey, volume)).To(Succeed(), "failed to get volume")
			g.Expect(volume.Spec.StoragePool.Name).To(Equal(storagePoolMatchingLabels.Name))
			g.Expect(volume.Status.State).To(Equal(storagev1alpha1.VolumeStatePending))
		}).Should(Succeed())
	})
})
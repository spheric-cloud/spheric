// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package etcd_test

import (
	"context"
	"io"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/apiserver/pkg/storage/etcd3/testserver"
	"k8s.io/utils/ptr"
	"spheric.cloud/spheric/actuo/codec"
	"spheric.cloud/spheric/actuo/list"
	"spheric.cloud/spheric/actuo/storage/etcd"
	"spheric.cloud/spheric/actuo/storage/store"
	"spheric.cloud/spheric/actuo/watch"
	"spheric.cloud/spheric/utils/generic"
)

type stringCodec struct{}

func (stringCodec) Encode(w io.Writer, obj *string) error {
	_, err := w.Write([]byte(*obj))
	return err
}

func (stringCodec) Decode(r io.Reader, into *string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	*into = string(data)
	return nil
}

var StringCodec codec.Codec[*string] = stringCodec{}

var _ = Describe("Etcd", func() {
	var (
		cl *clientv3.Client
		e  *etcd.Simple[*string, string]
	)

	BeforeEach(func() {
		cfg := testserver.NewTestConfig(GinkgoTB())
		cl = testserver.RunEtcd(GinkgoTB(), cfg)
		e = etcd.NewSimple[*string, string](
			cl,
			StringCodec,
			store.DefaultFactory[*string](),
			store.NopVersioner[*string, *list.List[*string, string]](),
		)
	})

	Describe("Create", func() {
		It("should create the object", func(ctx context.Context) {
			By("issuing a create")
			res, err := e.Create(ctx, "foo", generic.Pointer("obj"))
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(PointTo(Equal("obj")))

			By("getting the single kv from etcd")
			getResp, err := cl.Get(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Kvs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Key":   BeEquivalentTo("foo"),
				"Value": BeEquivalentTo("obj"),
			}))))
		})
	})

	Describe("Get", func() {
		It("should retrieve the object", func(ctx context.Context) {
			By("putting data into etcd")
			_, err := cl.Put(ctx, "foo", "bar")
			Expect(err).NotTo(HaveOccurred())

			By("retrieving the data")
			obj, err := e.Get(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(PointTo(Equal("bar")))
		})
	})

	Describe("Update", func() {
		It("should update the object", func(ctx context.Context) {
			By("putting data into etcd")
			_, err := cl.Put(ctx, "foo", "bar")
			Expect(err).NotTo(HaveOccurred())

			By("updating the data")
			obj, err := e.Update(ctx, "foo", false, func(ctx context.Context, oldObj *string) (newObj *string, err error) {
				Expect(*oldObj).To(Equal("bar"))
				return generic.Pointer("baz"), nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(PointTo(Equal("baz")))

			By("getting the single kv from etcd")
			getResp, err := cl.Get(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Kvs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Key":   BeEquivalentTo("foo"),
				"Value": BeEquivalentTo("baz"),
			}))))
		})
	})

	Describe("Delete", func() {
		It("should delete the object", func(ctx context.Context) {
			By("putting data into etcd")
			_, err := cl.Put(ctx, "foo", "bar")
			Expect(err).NotTo(HaveOccurred())

			By("deleting the object")
			obj, err := e.Delete(ctx, "foo", func(ctx context.Context, oldObj *string) error {
				Expect(*oldObj).To(Equal("bar"))
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(PointTo(Equal("bar")))

			By("verifying the kv is gone from etcd")
			getResp, err := cl.Get(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Kvs).To(BeEmpty())
		})
	})

	Describe("List", func() {
		It("should list all objects", func(ctx context.Context) {
			By("putting multiple objects under a key")
			_, err := cl.Put(ctx, "/root/o1", "foo")
			Expect(err).NotTo(HaveOccurred())

			_, err = cl.Put(ctx, "/root/o2", "bar")
			Expect(err).NotTo(HaveOccurred())

			By("listing all objects")
			objs, err := e.List(ctx, "/root")
			Expect(err).NotTo(HaveOccurred())
			Expect(objs).To(Equal(&list.List[*string, string]{
				Items: []string{"foo", "bar"},
			}))
		})
	})

	Describe("Watch", func() {
		It("should watch the object", func(ctx context.Context) {
			By("starting the watch")
			w, err := e.Watch(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			defer w.Stop()

			time.Sleep(2 * time.Second)

			By("creating an object")
			_, err = cl.Put(ctx, "foo", "bar")
			Expect(err).NotTo(HaveOccurred())

			By("waiting for the event")
			Eventually(ctx, w.Events()).Should(Receive(Equal(watch.Event[*string]{
				Type:   watch.EventTypeCreated,
				Object: ptr.To("bar"),
			})))
		})
	})
})

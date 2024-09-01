// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	"github.com/onsi/gomega/format"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/utils/generic"
	"spheric.cloud/spheric/utils/quota"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"spheric.cloud/spheric/utils/klog"
)

type DelegatingContext struct {
	lock sync.RWMutex
	ctx  context.Context
}

func (d *DelegatingContext) Deadline() (deadline time.Time, ok bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ctx.Deadline()
}

func (d *DelegatingContext) Done() <-chan struct{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ctx.Done()
}

func (d *DelegatingContext) Err() error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ctx.Err()
}

func (d *DelegatingContext) Value(key interface{}) interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ctx.Value(key)
}

func (d *DelegatingContext) Fulfill(ctx context.Context) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.ctx = ctx
}

func NewDelegatingContext(ctx context.Context) *DelegatingContext {
	return &DelegatingContext{ctx: ctx}
}

type ClientPromise struct {
	lock   sync.RWMutex
	client client.Client
}

func (d *ClientPromise) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Get(ctx, key, obj, opts...)
}

func (d *ClientPromise) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.List(ctx, list, opts...)
}

func (d *ClientPromise) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Create(ctx, obj, opts...)
}

func (d *ClientPromise) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Delete(ctx, obj, opts...)
}

func (d *ClientPromise) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Update(ctx, obj, opts...)
}

func (d *ClientPromise) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Patch(ctx, obj, patch, opts...)
}

func (d *ClientPromise) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.DeleteAllOf(ctx, obj, opts...)
}

func (d *ClientPromise) Status() client.SubResourceWriter {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Status()
}

func (d *ClientPromise) SubResource(subResource string) client.SubResourceClient {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.SubResource(subResource)
}

func (d *ClientPromise) Scheme() *runtime.Scheme {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.Scheme()
}

func (d *ClientPromise) RESTMapper() meta.RESTMapper {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.RESTMapper()
}

func (d *ClientPromise) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.GroupVersionKindFor(obj)
}

func (d *ClientPromise) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.client.IsObjectNamespaced(obj)
}

func (d *ClientPromise) FulfillWith(c client.Client, err error) error {
	if err != nil {
		return err
	}
	if c == nil {
		return fmt.Errorf("client is nil")
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	d.client = c
	return nil
}

func NewClientPromise() *ClientPromise {
	return &ClientPromise{client: defaultUninitializedClient}
}

var (
	errNotInitialized = errors.New("not initialized")

	emptyScheme                           = runtime.NewScheme()
	defaultUninitializedClient            = uninitializedClient{}
	defaultSubresourceUninitializedClient = subresourceUninitializedClient{}
	defaultUninitializedRESTMapper        = uninitializedRESTMapper{}
)

type uninitializedRESTMapper struct{}

func (uninitializedRESTMapper) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, errNotInitialized
}

func (uninitializedRESTMapper) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	return nil, errNotInitialized
}

func (uninitializedRESTMapper) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	return schema.GroupVersionResource{}, errNotInitialized
}

func (uninitializedRESTMapper) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	return nil, errNotInitialized
}

func (uninitializedRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	return nil, errNotInitialized
}

func (uninitializedRESTMapper) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	return nil, errNotInitialized
}

func (uninitializedRESTMapper) ResourceSingularizer(resource string) (singular string, err error) {
	return "", errNotInitialized
}

type subresourceUninitializedClient struct{}

func (subresourceUninitializedClient) Get(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceGetOption) error {
	return errNotInitialized
}

func (subresourceUninitializedClient) Create(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	return errNotInitialized
}

func (subresourceUninitializedClient) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	return errNotInitialized
}

func (subresourceUninitializedClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	return errNotInitialized
}

type uninitializedClient struct{}

func (uninitializedClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return errNotInitialized
}

func (uninitializedClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return errNotInitialized
}

func (uninitializedClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return errNotInitialized
}

func (uninitializedClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return errNotInitialized
}

func (uninitializedClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return errNotInitialized
}

func (uninitializedClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return errNotInitialized
}

func (uninitializedClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return errNotInitialized
}

func (uninitializedClient) Status() client.SubResourceWriter {
	return defaultSubresourceUninitializedClient
}

func (u uninitializedClient) SubResource(subResource string) client.SubResourceClient {
	return defaultSubresourceUninitializedClient
}

func (uninitializedClient) Scheme() *runtime.Scheme {
	return emptyScheme
}

func (uninitializedClient) RESTMapper() meta.RESTMapper {
	return defaultUninitializedRESTMapper
}

func (uninitializedClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, errNotInitialized
}

func (uninitializedClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return false, errNotInitialized
}

func SetupContext() context.Context {
	initCtx, initCancel := context.WithCancel(context.Background())
	delegCtx := NewDelegatingContext(initCtx)

	ginkgo.BeforeEach(func() {
		ctx, cancel := context.WithCancel(context.Background())
		ginkgo.DeferCleanup(cancel)

		delegCtx.Fulfill(ctx)
		if initCancel != nil {
			initCancel()
			initCancel = nil
		}
	})

	return delegCtx
}

// LowerCaseAlphabetCharset is a charset consisting of lower-case alphabet letters.
const LowerCaseAlphabetCharset = "abcdefghijklmnopqrstuvwxyz"

// RandomStringOptions are options for RandomString.
type RandomStringOptions struct {
	// Charset overrides the default RandomString charset if non-empty.
	Charset string
}

// ApplyToRandomString implements RandomStringOption.
func (o *RandomStringOptions) ApplyToRandomString(o2 *RandomStringOptions) {
	if o.Charset != "" {
		o2.Charset = o.Charset
	}
}

// ApplyOptions applies the slice of RandomStringOption to the RandomStringOptions.
func (o *RandomStringOptions) ApplyOptions(opts []RandomStringOption) {
	for _, opt := range opts {
		opt.ApplyToRandomString(o)
	}
}

// RandomStringOption is an option to RandomString.
type RandomStringOption interface {
	// ApplyToRandomString modifies the given RandomStringOptions with the option settings.
	ApplyToRandomString(o *RandomStringOptions)
}

// Charset specifies an explicit charset to use.
type Charset string

// ApplyToRandomString implements RandomStringOption.
func (s Charset) ApplyToRandomString(o *RandomStringOptions) {
	o.Charset = string(s)
}

// RandomString generates a random string of length n with the given options.
// If n is negative, RandomString panics.
func RandomString(n int, opts ...RandomStringOption) string {
	if n < 0 {
		panic("RandomString: negative length")
	}

	o := RandomStringOptions{}
	o.ApplyOptions(opts)

	charset := o.Charset
	if charset == "" {
		charset = LowerCaseAlphabetCharset
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteRune(rune(charset[rand.Intn(len(charset))]))
	}
	return sb.String()
}

// BeControlledBy matches any object that is controlled by the given owner.
func BeControlledBy(owner client.Object) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(obj client.Object) (bool, error) {
		return metav1.IsControlledBy(obj, owner), nil
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} be controlled by {{.Data 1}}", klog.KObjUID(owner))
}

type protoEqualMatcher struct {
	expected proto.Message
}

func ProtoEqual(expected proto.Message) types.GomegaMatcher {
	return &protoEqualMatcher{
		expected: expected,
	}
}

func (p *protoEqualMatcher) Match(actual interface{}) (bool, error) {
	actualMsg, ok := actual.(proto.Message)
	if !ok {
		return false, nil
	}

	return proto.Equal(actualMsg, p.expected), nil
}

func protoMessageToString(msg proto.Message) string {
	var sb strings.Builder
	formatProtoMessage(&sb, msg.ProtoReflect())
	return sb.String()
}

func formatProtoMessage(sb *strings.Builder, msg protoreflect.Message) {
	sb.WriteString(string(msg.Descriptor().FullName()))
	sb.WriteByte('(')

	var flds []protoreflect.FieldDescriptor
	msg.Range(func(fld protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		flds = append(flds, fld)
		return true
	})

	slices.SortFunc(flds, func(a, b protoreflect.FieldDescriptor) int { return int(int32(a.Number()) - int32(b.Number())) })

	for i, fld := range flds {
		if i > 0 {
			sb.WriteString(", ")
		}

		if fld.IsExtension() {
			// extension field: "[pkg.Msg.field]"
			sb.WriteString(string(fld.FullName()))
		} else if fld.Kind() != protoreflect.GroupKind {
			// ordinary field: "field"
			sb.WriteString(string(fld.Name()))
		} else {
			// group field: "MyGroup"
			//
			// The name of a group is the mangled version,
			// while the true name of a group is the message itself.
			// For example, for a group called "MyGroup",
			// the inlined message will be called "MyGroup",
			// but the field will be named "mygroup".
			// This rule complicates name logic everywhere.
			sb.WriteString(string(fld.Message().Name()))
		}
		sb.WriteString("=")
		formatProtoValue(sb, fld, msg.Get(fld))
	}
	sb.WriteRune(')')
}

func formatProtoValue(sb *strings.Builder, fld protoreflect.FieldDescriptor, v protoreflect.Value) {
	switch fld.Kind() {
	case protoreflect.GroupKind, protoreflect.MessageKind:
		formatProtoMessage(sb, v.Message())
	case protoreflect.EnumKind:
		// Invariant: only EnumValueDescriptor may appear here.
		desc := fld.Enum().Values().ByNumber(v.Enum())
		enum := fld.Parent()
		sb.WriteString(string(enum.Name()))
		sb.WriteRune('.')
		sb.WriteString(string(desc.Name()))
	default:
		_, _ = fmt.Fprint(sb, v.Interface())
	}
}

func (p *protoEqualMatcher) FailureMessage(actual interface{}) (message string) {
	expected := protoMessageToString(p.expected)

	actualR, ok := actual.(proto.Message)
	if !ok {
		return fmt.Sprintf("Expected\n%s\n%s\n%s", format.Object(actual, 1), "to equal", format.IndentString(expected, 1))
	}

	actualS := protoMessageToString(actualR)
	return fmt.Sprintf("Expected\n%s\n%s\n%s", format.IndentString(actualS, 1), "to equal", format.IndentString(expected, 1))
}

func (p *protoEqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", p.expected)
}

type equalResourcesMatcher struct {
	expected corev1alpha1.ResourceList
}

func EqualResources(expected corev1alpha1.ResourceList) types.GomegaMatcher {
	return &equalResourcesMatcher{
		expected: expected,
	}
}

func (e *equalResourcesMatcher) Match(actual interface{}) (bool, error) {
	actualR, err := generic.Cast[corev1alpha1.ResourceList](actual)
	if err != nil {
		return false, err
	}

	return quota.Equals(e.expected, actualR), nil
}

func formatResources(res corev1alpha1.ResourceList) string {
	keys := slices.Collect(maps.Keys(res))
	slices.Sort(keys)

	var sb strings.Builder
	sb.WriteRune('{')
	for i, key := range keys {
		if i != 0 {
			sb.WriteString(", ")
		}

		v := res[key]
		sb.WriteString(string(key))
		sb.WriteRune(':')
		sb.WriteString(v.String())
	}
	sb.WriteRune('}')
	return sb.String()
}

func (e *equalResourcesMatcher) FailureMessage(actual interface{}) string {
	expected := formatResources(e.expected)

	actualR, ok := actual.(corev1alpha1.ResourceList)
	if !ok {
		return fmt.Sprintf("Expected\n%s\n%s\n%s", format.Object(actual, 1), "to equal", format.IndentString(expected, 1))
	}

	actualS := formatResources(actualR)
	return fmt.Sprintf("Expected\n%s\n%s\n%s", format.IndentString(actualS, 1), "to equal", format.IndentString(expected, 1))
}

func (e *equalResourcesMatcher) NegatedFailureMessage(actual interface{}) string {
	actualR, ok := actual.(corev1alpha1.ResourceList)
	if !ok {
		actual = formatResources(actualR)
	}

	expected := formatResources(e.expected)
	return format.Message(actual, "not to equal", expected)
}

// SetupNamespace sets up a namespace before each test and tears the namespace down after each test.
func SetupNamespace(c client.Client) *corev1.Namespace {
	return SetupNewObject[*corev1.Namespace](c, func(ns *corev1.Namespace) {
		*ns = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-ns",
			},
		}
	})
}

func SetupObject(c client.Client, obj client.Object, f func()) {
	ginkgo.BeforeEach(func(ctx context.Context) {
		f()
		gomega.Expect(c.Create(ctx, obj)).To(gomega.Succeed(), "failed to create object %T (%s)", obj, client.ObjectKeyFromObject(obj))
		ginkgo.DeferCleanup(DeleteIgnoreNotFound(c, obj))
	})
}

func SetupNewObject[O interface {
	client.Object
	*OStruct
}, OStruct any](c client.Client, f func(obj O)) O {
	obj := O(new(OStruct))
	SetupObject(c, obj, func() {
		f(obj)
	})
	return obj
}

// DeleteIgnoreNotFound returns a function to clean up an object if it exists.
func DeleteIgnoreNotFound(c client.Client, obj client.Object) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := c.Delete(ctx, obj)
		return client.IgnoreNotFound(err)
	}
}

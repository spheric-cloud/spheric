// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ConfigMapKeySelector is a reference to a specific 'key' within a ConfigMap resource.
// In some instances, `key` is a required field.
// +structType=atomic
type ConfigMapKeySelector struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name,omitempty"`
	// The key of the entry in the ConfigMap resource's `data` field to be used.
	// Some instances of this field may be defaulted, in others it may be
	// required.
	// +optional
	Key string `json:"key,omitempty"`
}

// SecretKeySelector is a reference to a specific 'key' within a Secret resource.
// In some instances, `key` is a required field.
// +structType=atomic
type SecretKeySelector struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name,omitempty"`
	// The key of the entry in the Secret resource's `data` field to be used.
	// Some instances of this field may be defaulted, in others it may be
	// required.
	// +optional
	Key string `json:"key,omitempty"`
}

// ObjectSelector specifies how to select objects of a certain kind.
type ObjectSelector struct {
	// Kind is the kind of object to select.
	Kind string `json:"kind"`
	// LabelSelector is the label selector to select objects of the specified Kind by.
	metav1.LabelSelector `json:",inline"`
}

// LocalUIDReference is a reference to another entity including its UID
// +structType=atomic
type LocalUIDReference struct {
	// Name is the name of the referenced entity.
	Name string `json:"name"`
	// UID is the UID of the referenced entity.
	UID types.UID `json:"uid"`
}

func LocalObjUIDRef(obj metav1.Object) LocalUIDReference {
	return LocalUIDReference{
		Name: obj.GetName(),
		UID:  obj.GetUID(),
	}
}

func NewLocalObjUIDRef(obj metav1.Object) *LocalUIDReference {
	return &LocalUIDReference{
		Name: obj.GetName(),
		UID:  obj.GetUID(),
	}
}

// UIDReference is a reference to another entity in a potentially different namespace including its UID.
// +structType=atomic
type UIDReference struct {
	// Namespace is the namespace of the referenced entity. If empty,
	// the same namespace as the referring resource is implied.
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the referenced entity.
	Name string `json:"name"`
	// UID is the UID of the referenced entity.
	UID types.UID `json:"uid,omitempty"`
}

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
type LocalObjectReference struct {
	// Name of the referent.
	Name string `json:"name,omitempty"`
}

func LocalObjRef(name string) LocalObjectReference {
	return LocalObjectReference{Name: name}
}

func NewLocalObjRef(name string) *LocalObjectReference {
	return &LocalObjectReference{Name: name}
}

// Taint marks an effect with a value on a target resource pool.
type Taint struct {
	// The taint key to be applied to a resource pool.
	Key string `json:"key"`
	// The taint value corresponding to the taint key.
	Value string `json:"value,omitempty"`
	// The effect of the taint on resources
	// that do not tolerate the taint.
	// Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
	Effect TaintEffect `json:"effect"`
}

type TaintEffect string

const (
	// TaintEffectNoSchedule causes not to allow new resources to schedule onto the resource pool unless they tolerate
	// the taint, but allow all already-running resources to continue running.
	// Enforced by the scheduler.
	TaintEffectNoSchedule TaintEffect = "NoSchedule"
)

// Toleration marks the resource the toleration is attached to tolerate any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	Key string `json:"key,omitempty"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a resource can
	// tolerate all taints of a particular category.
	Operator TolerationOperator `json:"operator,omitempty"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	Value string `json:"value,omitempty"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule.
	Effect TaintEffect `json:"effect,omitempty"`
}

// ToleratesTaint checks if the toleration tolerates the taint.
// The matching follows the rules below:
// (1) Empty toleration.effect means to match all taint effects,
//
//	otherwise taint effect must equal to toleration.effect.
//
// (2) If toleration.operator is 'Exists', it means to match all taint values.
// (3) Empty toleration.key means to match all taint keys.
//
//	If toleration.key is empty, toleration.operator must be 'Exists';
//	this combination means to match all taint values and all taint keys.
func (t *Toleration) ToleratesTaint(taint *Taint) bool {
	if len(t.Effect) > 0 && t.Effect != taint.Effect {
		return false
	}

	if len(t.Key) > 0 && t.Key != taint.Key {
		return false
	}

	switch t.Operator {
	case "", TolerationOpEqual: // empty operator means Equal
		return t.Value == taint.Value
	case TolerationOpExists:
		return true
	default:
		return false
	}
}

// TolerationOperator is the set of operators that can be used in a toleration.
type TolerationOperator string

const (
	TolerationOpEqual  TolerationOperator = "Equal"
	TolerationOpExists TolerationOperator = "Exists"
)

// TolerateTaints returns if tolerations tolerate all taints
func TolerateTaints(tolerations []Toleration, taints []Taint) bool {
Outer:
	for _, taint := range taints {
		for _, toleration := range tolerations {
			if toleration.ToleratesTaint(&taint) {
				continue Outer
			}
		}
		return false
	}
	return true
}

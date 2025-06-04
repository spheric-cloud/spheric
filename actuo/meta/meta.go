// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	pt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = pt.Local()
	return nil
}

type Object interface {
	GetNamespace() string
	SetNamespace(namespace string)
	GetName() string
	SetName(name string)
	GetResourceVersion() string
	SetResourceVersion(resourceVersion string)
	GetGeneration() int64
	SetGeneration(generation int64)
	GetCreationTimestamp() Time
	SetCreationTimestamp(t Time)
	GetDeletionTimestamp() *Time
	SetDeletionTimestamp(t *Time)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
}

type ObjectMeta struct {
	Namespace         string            `json:"namespace,omitempty"`
	Name              string            `json:"name,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	Generation        int64             `json:"generation,omitempty"`
	CreationTimestamp Time              `json:"creationTimestamp,omitempty"`
	DeletionTimestamp *Time             `json:"deletionTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Finalizers        []string          `json:"finalizers,omitempty"`
}

func (o *ObjectMeta) GetNamespace() string                         { return o.Namespace }
func (o *ObjectMeta) SetNamespace(namespace string)                { o.Namespace = namespace }
func (o *ObjectMeta) GetName() string                              { return o.Name }
func (o *ObjectMeta) SetName(name string)                          { o.Name = name }
func (o *ObjectMeta) GetResourceVersion() string                   { return o.ResourceVersion }
func (o *ObjectMeta) SetResourceVersion(resourceVersion string)    { o.ResourceVersion = resourceVersion }
func (o *ObjectMeta) GetGeneration() int64                         { return o.Generation }
func (o *ObjectMeta) SetGeneration(generation int64)               { o.Generation = generation }
func (o *ObjectMeta) GetCreationTimestamp() Time                   { return o.CreationTimestamp }
func (o *ObjectMeta) SetCreationTimestamp(t Time)                  { o.CreationTimestamp = t }
func (o *ObjectMeta) GetDeletionTimestamp() *Time                  { return o.DeletionTimestamp }
func (o *ObjectMeta) SetDeletionTimestamp(t *Time)                 { o.DeletionTimestamp = t }
func (o *ObjectMeta) GetLabels() map[string]string                 { return o.Labels }
func (o *ObjectMeta) SetLabels(labels map[string]string)           { o.Labels = labels }
func (o *ObjectMeta) GetAnnotations() map[string]string            { return o.Annotations }
func (o *ObjectMeta) SetAnnotations(annotations map[string]string) { o.Annotations = annotations }
func (o *ObjectMeta) GetFinalizers() []string                      { return o.Finalizers }

func (o *ObjectMeta) SetFinalizers(finalizers []string) { o.Finalizers = finalizers }

type List interface {
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetContinue() string
	SetContinue(c string)
}

type ListMeta struct {
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty"`
}

func (m *ListMeta) GetResourceVersion() string                { return m.ResourceVersion }
func (m *ListMeta) SetResourceVersion(resourceVersion string) { m.ResourceVersion = resourceVersion }
func (m *ListMeta) GetContinue() string                       { return m.Continue }
func (m *ListMeta) SetContinue(c string)                      { m.Continue = c }

type WatchEvent struct {
	Type   string `json:"type"`
	Object any    `json:"object"`
}

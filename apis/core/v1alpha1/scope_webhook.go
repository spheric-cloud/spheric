/*
 * Copyright (c) 2021 by the OnMetal authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var scopelog = logf.Log.WithName("scope-resource")

//+kubebuilder:webhook:path=/mutate-core-onmetal-de-v1alpha1-scope,mutating=true,failurePolicy=fail,sideEffects=None,groups=core.onmetal.de,resources=scopes,verbs=create;update,versions=v1alpha1,name=mscope.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Scope{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Scope) Default() {
	scopelog.Info("default", "name", r.Name)

	if r.Status.State == "" {
		r.Status.State = ScopeStateInitial
	}
}

//+kubebuilder:webhook:path=/validate-core-onmetal-de-v1alpha1-scope,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.onmetal.de,resources=scopes,verbs=create;update;delete,versions=v1alpha1,name=vscope.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Scope{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Scope) ValidateCreate() error {
	scopelog.Info("validate create", "name", r.Name)
	return r.validateScope()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Scope) ValidateUpdate(old runtime.Object) error {
	scopelog.Info("validate update", "name", r.Name)
	return r.validateScopeUpdate(old)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Scope) ValidateDelete() error {
	scopelog.Info("validate delete", "name", r.Name)
	return r.validateScopeDelete()
}
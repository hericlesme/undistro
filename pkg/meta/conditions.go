/*
Copyright 2020 The UnDistro authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package meta

import (
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ReadyCondition is the name of the Ready condition implemented by all toolkit
	// resources.
	ReadyCondition string = "Ready"
)

const (
	// ReconciliationSucceededReason represents the fact that the reconciliation of
	// a toolkit resource has succeeded.
	ReconciliationSucceededReason string = "ReconciliationSucceeded"

	// ReconciliationFailedReason represents the fact that the reconciliation of a
	// toolkit resource has failed.
	ReconciliationFailedReason string = "ReconciliationFailed"

	// ProgressingReason represents the fact that the reconciliation of a toolkit
	// resource is underway.
	ProgressingReason string = "Progressing"

	// DependencyNotReadyReason represents the fact that one of the toolkit resource
	// dependencies is not ready.
	DependencyNotReadyReason string = "DependencyNotReady"

	// PausedReason represents the fact that the reconciliation of a toolkit
	// resource is Paused.
	PausedReason string = "Paused"

	// URLInvalidReason represents the fact that a given source has an invalid URL.
	URLInvalidReason string = "URLInvalid"

	// StorageOperationFailedReason signals a failure caused by a storage operation.
	StorageOperationFailedReason string = "StorageOperationFailed"

	// AuthenticationFailedReason represents the fact that a given secret does not
	// have the required fields or the provided credentials do not match.
	AuthenticationFailedReason string = "AuthenticationFailed"

	// VerificationFailedReason represents the fact that the cryptographic
	// provenance verification for the source failed.
	VerificationFailedReason string = "VerificationFailed"
)

// InReadyCondition returns if the given Condition slice has a ReadyCondition
// with a 'True' condition status.
func InReadyCondition(conditions []metav1.Condition) bool {
	return apimeta.IsStatusConditionTrue(conditions, ReadyCondition)
}

// ObjectWithStatusConditions is an interface that describes kubernetes resource
// type structs with Status Conditions
type ObjectWithStatusConditions interface {
	GetStatusConditions() *[]metav1.Condition
}

// SetResourceCondition sets the given condition with the given status,
// reason and message on a resource.
func SetResourceCondition(obj ObjectWithStatusConditions, condition string, status metav1.ConditionStatus, reason, message string) {
	conditions := obj.GetStatusConditions()

	newCondition := metav1.Condition{
		Type:    condition,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	apimeta.SetStatusCondition(conditions, newCondition)
}

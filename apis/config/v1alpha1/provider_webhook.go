/*


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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var providerlog = logf.Log.WithName("provider-resource")

func (r *Provider) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-config-undistro-io-v1alpha1-provider,mutating=true,failurePolicy=fail,groups=config.undistro.io,resources=providers,verbs=create;update;delete,versions=v1alpha1,name=mprovider.undistro.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Defaulter = &Provider{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Provider) Default() {
	providerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-config-undistro-io-v1alpha1-provider,mutating=false,failurePolicy=fail,groups=config.undistro.io,resources=providers,versions=v1alpha1,name=vprovider.undistro.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Validator = &Provider{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Provider) ValidateCreate() error {
	providerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Provider) ValidateUpdate(old runtime.Object) error {
	providerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Provider) ValidateDelete() error {
	providerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package v1alpha1

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	defaultKubernetesVersion          = "v1.18.2"
	defaultControlPlaneReplicas int64 = 1
	defaultWorkerReplicas       int64 = 0
)

// log is for logging in this package.
var clusterlog = logf.Log.WithName("cluster-resource")

func (r *Cluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-getupcloud-com-v1alpha1-cluster,mutating=true,failurePolicy=fail,groups=getupcloud.com,resources=clusters,verbs=create;update,versions=v1alpha1,name=mcluster.getupcloud.com

var _ webhook.Defaulter = &Cluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Cluster) Default() {
	clusterlog.Info("default", "name", r.Name)
	if r.Spec.CniName == EmptyCNI {
		r.Spec.CniName = CalicoCNI
	}
	if r.Spec.KubernetesVersion == "" {
		r.Spec.KubernetesVersion = defaultKubernetesVersion
	}
	if r.Spec.ControlPlaneNode.Replicas == nil {
		r.Spec.ControlPlaneNode.Replicas = new(int64)
		*r.Spec.ControlPlaneNode.Replicas = defaultControlPlaneReplicas
	}
	if r.Spec.WorkerNode.Replicas == nil {
		r.Spec.WorkerNode.Replicas = new(int64)
		*r.Spec.WorkerNode.Replicas = defaultWorkerReplicas
	}
}
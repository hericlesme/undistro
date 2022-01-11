package app

const (
	KubernetesVersionChanged        = "Kubernetes version changed"
	ControlPlaneReplicasChanged     = "Number of control plane replicas changed"
	ControlPlaneMachineTypeChanged  = "Machine type changed in control plane"
	ControlPlaneLabelsChanged       = "Control plane labels have changed"
	ControlPlaneTaintsChanged       = "Taints have changed on the control plane"
	ControlPlaneProviderTagsChanged = "Infrastructure provider tags have changed in the control plane"
	WorkersReplicasChanged          = "Number of worker replicas changed"
	WorkersChanged                  = "Number of workers changed"
	WorkerMachineTypeChanged        = "Machine type changed in worker pool"
	WorkerLabelsChanged             = "Worker labels have change"
	WorkerTaintsChanged             = "Taints have changed on the worker"
	WorkerProviderTagsChanged       = "Infrastructure provider tags have changed in the worker"
	WorkerAutoscalingChanged        = "Worker autoscaling have change"
)

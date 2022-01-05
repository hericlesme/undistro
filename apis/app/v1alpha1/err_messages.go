package v1alpha1

const (
	SshKeyRequired          = "The 'sshKey' field is required when bastion is enabled"
	SshRequiredInEC2        = "The 'sshKey' field is required when flavor is ec2"
	FlavorRequired          = "The 'flavor' field must to be populated"
	FlavorNotValid          = "The target flavor is not valid"
	CPRequiredInNonManaged  = "The 'controlPlane' field is required in non-managed clusters"
	CPRequiredInSelfHosted  = "The 'controlPlane' field is required when is a self hosted cluster"
	InvalidSemVer           = "The 'kubernetesVersion' field must to be a semantic versioning"
	UpdateClusterNotReady   = "Can't update cluster that isn't ready"
	ImmutableField          = "The target field is immutable"
	NetAddrConflict         = "ID or CIDRBlock must be set to avoid network conflicts with others clusters"
	InvalidClusterNameInAws = "Invalid cluster name for AWS" // Add valid ones
)

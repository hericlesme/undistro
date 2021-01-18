export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=AKIAJ2HNRPZEWXY4UVKA
export AWS_SECRET_ACCESS_KEY='h1kyrDLt+z4itVFu7g1MEM+ZOcgE/5PfdpC+9MwR'
export EXP_EKS=true
export EXP_EKS_IAM=true
export EXP_EKS_ADD_ROLES=true
export EXP_CLUSTER_RESOURCE_SET=true
export EXP_MACHINE_POOL=true
export AWS_REGION=us-east-1
export AWS_SSH_KEY_NAME=undistro
# Select instance types
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
clusterawsadm bootstrap iam create-cloudformation-stack --config bootstrap-config.yaml
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

# Finally, initialize the management cluster
clusterctl init --infrastructure=aws --control-plane aws-eks --bootstrap aws-eks

clusterctl config cluster capi-eks-quickstart --flavor eks-managedmachinepool --kubernetes-version v1.18.9 --worker-machine-count=3 > demos/capi-eks-quickstart.yaml

kubectl apply -f demos/capi-eks-quickstart.yaml

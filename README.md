# UnDistro Kubernetes platform

## Documentation

### Requirements

- Install and setup [kubectl] in your local environment
- Install [Kind] and [Docker]

### Install and/or configure a kubernetes cluster

UnDistro requires an existing Kubernetes cluster accessible via kubectl; during the installation process the
Kubernetes cluster will be transformed into a [management cluster] by installing the UnDistro [provider components], so it
is recommended to keep it separated from any application workload.

It is a common practice to create a temporary, local bootstrap cluster which is then used to provision
a target [management cluster] on the selected [infrastructure provider].

Choose one of the options below:

1. **Existing Management Cluster**

For production use-cases a "real" Kubernetes cluster should be used with appropriate backup and DR policies and procedures in place.

```bash
export KUBECONFIG=<...>
```

2. **Kind**

[kind] is not designed for production use.

**Minimum [kind] supported version**: v0.9.0

</aside>

[kind] can be used for creating a local Kubernetes cluster for development environments or for
the creation of a temporary [bootstrap cluster] used to provision a target [management cluster] on the selected infrastructure provider.

### Quick start for AWS

**Create configuration file**
replace `<key>` to your keys

```yaml
providers:
  -
    name: aws
    configuration:
      accessKeyID: <key>
      secretAccessKey: <key>
      sessionToken: <key> # if you use 2FA
```

### Install UnDistro CLI
The UnDistro CLI tool handles the lifecycle of a UnDistro management cluster.

Download the latest release from releases page https://github.com/getupio-undistro/undistro/releases.

### Initialize the management cluster

Now that we've got UnDistro CLI installed and all the prerequisites in place, let's transform the Kubernetes cluster
into a management cluster by using `undistro install`.

```bash
undistro --config undistro-config.yaml install   
```

### Adding SSH key on AWS account

Follow the AWS documentation https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html

### Creating a self hosted kubernetes on AWS

We will create a cluster with 3 controlplane node and 3 worker nodes

```yaml
apiVersion: app.undistro.io/v1alpha1
kind: Cluster
metadata:
  name: undistro-quickstart
  namespace: default
spec:
  kubernetesVersion: v1.19.5
  controlPlane:
    replicas: 3
    machineType: t3.medium
  workers:
    - replicas: 3
      machineType: t3.medium
  infrastructureProvider:
    name: aws
    sshKey: <YOUR SSH KEY NAME>
    region: us-east-1
```

### Creating an EKS kubernetes on AWS

We will create a cluster with 3 controlplane node and 3 worker nodes

```yaml
apiVersion: app.undistro.io/v1alpha1
kind: Cluster
metadata:
  name: undistro-quickstart
  namespace: default
spec:
  kubernetesVersion: v1.19.5
  workers:
    - replicas: 3
      machineType: t3.medium
  infrastructureProvider:
    name: aws
    sshKey: <YOUR SSH KEY NAME>
    region: us-east-1
    managed: true
```
**create a file with content above.**

```
undistro create -f cluster-aws.yaml
```

### Waiting cluster to be ready to use

When cluster is ready you will see **Cluster reconciliation succeeded** in command below output:

```
undistro get clusters
```

### Getting cluster kubeconfig

```
undistro get kubeconfig undistro-quickstart -n default
```

### Making a remote cluster to be your management cluster

```
undistro --config undistro-config.yaml move undistro-quickstart -n default
```

[Docker]: https://www.docker.com/
[kind]: https://kind.sigs.k8s.io/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
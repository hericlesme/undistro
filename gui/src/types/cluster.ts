import type { KubernetesObject } from '@kubernetes/client-node'

export interface Cluster {
  name: string
  provider: string
  flavor: string
  k8sVersion: string
  clusterGroup: string
  machines: number
  age: string | number // based on timestamp
  status: string
}

export interface Worker {
  id: string
  name: string
  machineType: string
  replicas: number
  infraNode: boolean
}

export interface KubernetesResource extends KubernetesObject {
  spec: {
    clusterName: string
  }
}
export interface Subnet {
  id: string
  cidrBlock: string
  zone: string
  isPublic: boolean
}

export interface Network {
  vpc: {
    id: string
    cidrBlock: string
  }
}

export interface InfrastructureProvider {
  flavor: string
  name: string
  region: string
  sshKey: string
}

export interface ControlPlane {
  machineType: string
  replicas: number
}

export interface CreateClusterRequest extends KubernetesObject {
  spec: {
    kubernetesVersion: string
    controlPlane: ControlPlane
    infrastructureProvider: InfrastructureProvider
    workers: Worker[]
    network: Network
  }
}

export interface Provider extends KubernetesObject {
  metadata: {
    name: string
  }
  status: {
    regionNames?: string[]
  }
}

export interface Flavor {
  name: string
  supportedVersions: string[]
}

export interface MachineType {
  name: string
  mem: number
  cpu: number
  zones: string[]
}

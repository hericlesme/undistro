import type { Cluster } from '@/types/cluster'
import { getAge, getStatusFromConditions } from '@/lib/cluster'

export function getMetadataUrl(kind: string): string {
  return `apis/metadata.undistro.io/v1alpha1/${kind}`
}

export function clusterDataHandler(clusterList): Cluster[] {
  return clusterList.items.map(cl => {
    let machines = 0
    let workers = cl.spec.workers as Array<any>

    if (workers) {
      workers.forEach(w => {
        machines += w.replicas as number
      })
    }

    let conditions = cl.status.conditions
    if (cl.spec.controlPlane != undefined && cl.spec.controlPlane.replicas != undefined) {
      machines += cl.spec.controlPlane.replicas as number
    }

    return {
      name: cl.metadata.name as string,
      provider: cl.spec.infrastructureProvider.name as string,
      flavor: cl.spec.infrastructureProvider.flavor as string,
      k8sVersion: cl.spec.kubernetesVersion as string,
      clusterGroup: cl.metadata.namespace as string,
      machines: machines,
      age: getAge(cl.metadata.creationTimestamp as string),
      status: getStatusFromConditions(conditions),
      workers: workers || [],
      controlPlane: cl.spec.controlPlane as any
    }
  })
}

export function machineTypeDataHandler(machines) {
  return machines.map(machine => ({
    zones: machine.spec.availabilityZones,
    name: machine.metadata.name,
    mem: machine.spec.memory,
    cpu: machine.spec.vcpus
  }))
}

export function providersDataHandler(providers) {
  return providers.filter((provider: any) => {
    return provider.spec.category.includes('infra')
  })
}

export function flavorsDataHandler(flavors) {
  return flavors.map(flavor => ({
    name: flavor.metadata.name,
    supportedVersions: flavor.spec.supportedK8SVersions
  }))
}

import { Cluster, getAge, getStatusFromConditions } from '@/lib/cluster'

export function clusterDataHandler(clusterList): Cluster[] {
  return clusterList.items.map(cl => {
    let machines = 0
    let workers = cl.spec.workers as Array<any>
    workers.forEach(w => {
      machines += w.replicas as number
    })
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
      status: getStatusFromConditions(conditions)
    }
  })
}

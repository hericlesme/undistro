/* eslint-disable react-hooks/exhaustive-deps */
import { useEffect, useState } from 'react'
import Table from '@components/nodepoolTable'
import Api from 'util/api'
import moment from 'moment'

const headers = [
  { name: 'Name', field: 'name' },
  { name: 'Type', field: 'type' },
  { name: 'Replicas', field: 'replicas' },
  { name: 'K8s Version', field: 'version' },
  { name: 'Labels', field: 'labels' },
  { name: 'Taints', field: 'taints' },
  { name: 'Age', field: 'age' },
  { name: 'Status', field: 'status' }
]

export default function NodepoolsPage() {
  const [nodePools, setNodePools] = useState<any>([])

  const getControlPlanesAndWorkers = async () => {
    const clusters = await Api.Cluster.list()

    const controlPlanes = clusters.items.map((cluster: any) => {
      const message: string = cluster.status.conditions[0].message.toLowerCase()

      return {
        ...cluster.spec.controlPlane,
        age: moment(cluster.metadata.creationTimestamp)
          .startOf('day')
          .fromNow(),
        labels: 0,
        name: cluster.metadata.name,
        status: getStatus(message),
        taints: 0,
        type: 'Control Plane',
        version: cluster.spec.kubernetesVersion
      }
    })

    const workers = clusters.items.reduce((workers: any, cluster: any) => {
      const message: string = cluster.status.conditions[0].message.toLowerCase()

      return [
        ...workers,
        ...cluster.spec.workers.map((worker: any, i: number) => ({
          ...worker,
          age: moment(cluster.metadata.creationTimestamp)
            .startOf('day')
            .fromNow(),
          labels: Object.keys(worker.labels || {}).length,
          name: `${cluster.metadata.name}-mp-${i}`,
          status: getStatus(message),
          taints: (worker.taints || []).length,
          type: worker.infraNode ? 'InfraNode' : 'Worker',
          version: cluster.spec.kubernetesVersion
        }))
      ]
    }, [])

    return [...controlPlanes, ...workers]
  }

  const getStatus = (message: string) => {
    if (message.includes('wait cluster')) return 'Provisioning'
    else if (message.includes('error')) return 'Error'
    else if (message.includes('paused')) return 'Paused'
    else if (message.includes('deleting')) return 'Deleting'
    else return 'Ready'
  }

  const getNodePools = async () => {
    try {
      const controlPlanesAndWorkers = await getControlPlanesAndWorkers()

      const nodePools = controlPlanesAndWorkers.map((controlPlaneOrWorker: any) => {
        return {
          name: controlPlaneOrWorker.name,
          type: controlPlaneOrWorker.type,
          replicas: controlPlaneOrWorker.replicas,
          version: controlPlaneOrWorker.version,
          labels: controlPlaneOrWorker.labels,
          taints: controlPlaneOrWorker.taints,
          age: controlPlaneOrWorker.age,
          status: controlPlaneOrWorker.status
        }
      })

      setNodePools(nodePools)
    } catch (err) {
      console.log(err)
    }
  }

  useEffect(() => {
    getNodePools()
  }, [])

  return (
    <div className="home-page-route">
      <Table data={nodePools || []} header={headers} />
    </div>
  )
}

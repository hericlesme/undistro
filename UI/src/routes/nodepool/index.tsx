/* eslint-disable react-hooks/exhaustive-deps */
import { useEffect, useState } from 'react'
import Table from '@components/nodepoolTable'
import moment from 'moment'
import { useClusters } from 'providers/ClustersProvider'
import BreadCrumb from '@components/breadcrumb'
import { useServices } from 'providers/ServicesProvider'
import { Layout } from '@components/layout'

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
  const { Api } = useServices()
  const [nodePools, setNodePools] = useState<any>([])
  const { clusters } = useClusters()

  const getControlPlanesAndWorkers = async () => {
    let allControlPlanes: any = []
    let allWorkers: any = []

    await Promise.all(
      clusters.map(async cluster => {
        const fetchedCluster = await Api.Cluster.get(cluster.namespace, cluster.name)
        const message: string = fetchedCluster.status.conditions[0].message.toLowerCase()
        const controlPlane = {
          ...fetchedCluster.spec.controlPlane,
          age: moment(fetchedCluster.metadata.creationTimestamp)
            .startOf('day')
            .fromNow(),
          labels: 0,
          name: fetchedCluster.metadata.name,
          status: getStatus(message),
          taints: 0,
          type: 'Control Plane',
          version: fetchedCluster.spec.kubernetesVersion
        }

        const workers = fetchedCluster.spec.workers.map((worker: any, i: number) => ({
          ...worker,
          age: moment(fetchedCluster.metadata.creationTimestamp)
            .startOf('day')
            .fromNow(),
          labels: Object.keys(worker.labels || {}).length,
          name: `${fetchedCluster.metadata.name}-mp-${i}`,
          status: getStatus(message),
          taints: (worker.taints || []).length,
          type: worker.infraNode ? 'InfraNode' : 'Worker',
          version: fetchedCluster.spec.kubernetesVersion
        }))

        allControlPlanes.push(controlPlane)
        allWorkers = [...allWorkers, ...workers]
      })
    )

    return [...allControlPlanes, ...allWorkers]
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

  const getClusterName = () => {
    return clusters.length > 1 ? 'Multiple Clusters' : clusters.map(elm => elm.name)
  }

  const routes = [
    { name: 'Clusters', url: '/' },
    { name: getClusterName(), url: '/' }
  ]

  useEffect(() => {
    getNodePools()
  }, [])

  return (
    <Layout>
      <div className="home-page-route">
        <BreadCrumb routes={routes} />
        <Table data={nodePools || []} header={headers} />
      </div>
    </Layout>
  )
}

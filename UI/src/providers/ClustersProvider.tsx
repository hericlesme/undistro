import { createContext, useState, useContext } from 'react'
import type { ReactNode } from 'react'

type Cluster = {
  name: string
  namespace: string
  paused: boolean
}

type ClustersContextValue = {
  clusters: Cluster[]
  addCluster: (cluster: Cluster) => void
  addClusters: (clusters: Cluster[]) => void
  removeCluster: (name: string) => void
  clear: () => void
  isEmpty: boolean
}

const ClustersContext = createContext({} as ClustersContextValue)

export const ClustersProvider = ({ children } : { children: ReactNode }) => {
  const [clusters, setClusters] = useState<Cluster[]>([])

  const addCluster = (cluster: Cluster) => {
    setClusters([...clusters, cluster])
  }

  const addClusters = (clusters: Cluster[]) => {
    setClusters(clusters)
  }

  const removeCluster = (name: string) => {
    setClusters([...clusters.filter(cluster => cluster.name !== name)])
  }

  const clear = () => {
    setClusters([])
  }

  return <ClustersContext.Provider value={{ clusters, addCluster, addClusters, removeCluster, clear, isEmpty: clusters.length === 0 }}>{children}</ClustersContext.Provider>
}

export const useClusters = () => {
  const { clusters, addCluster, removeCluster, clear, addClusters, isEmpty } = useContext(ClustersContext)

  return {
    clusters,
    addCluster,
    removeCluster,
    clear,
    addClusters,
    isEmpty
  }
}

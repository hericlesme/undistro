import { createContext, useContext, useEffect, useState } from 'react'

export type ClusterContextType = {
  clusters: string[]
  setClusters: (clusters: string[]) => void
}

export const ClusterContext = createContext<ClusterContextType>({ clusters: [], setClusters: clusters => {} })
export const useClusters = () => useContext(ClusterContext)

export const ClusterProvider: React.FC<{}> = ({ children }) => {
  const [clusters, setClusters] = useState<string[]>([])

  return <ClusterContext.Provider value={{ clusters, setClusters }}>{children}</ClusterContext.Provider>
}

import type { FC } from 'react'
import { useState } from 'react'
import { Topbar } from '@/components/Topbar'
import { Workarea } from '@/components/Workarea'
import { ClusterContext } from '@/contexts/ClusterContext'

type WorkspaceProps = {
  children?: React.ReactNode
  selectedClusters: string[]
}

const Workspace: FC<WorkspaceProps> = (props: WorkspaceProps) => {
  const [clusters, setClusters] = useState<string[]>(props.selectedClusters)
  return (
    <>
      <ClusterContext.Provider value={{ clusters, setClusters }}>
        <Topbar />
        <Workarea>{props.children}</Workarea>
      </ClusterContext.Provider>
    </>
  )
}

export { Workspace }

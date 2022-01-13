import type { FC } from 'react'
import { Topbar } from '@/components/Topbar'
import { Workarea } from '@/components/Workarea'

type WorkspaceProps = {
  children?: React.ReactNode
  selectedClusters: string[]
}

const Workspace: FC<WorkspaceProps> = (props: WorkspaceProps) => {
  return (
    <>
      <Topbar />
      <Workarea>{props.children}</Workarea>
    </>
  )
}

export { Workspace }

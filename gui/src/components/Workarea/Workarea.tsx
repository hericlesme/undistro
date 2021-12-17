import type { ReactNode, FC } from 'react'
import { LeftMenuArea } from '@/components/LeftMenu'
import styles from '@/components/Workarea/Workarea.module.css'

type WorkareaProps = {
  children?: ReactNode
}

const Workarea: FC<WorkareaProps> = (props: WorkareaProps) => {
  return (
    <>
      <div className={styles.mainWorkspaceArea}>
        <div className={styles.leftMenuArea}>
          <LeftMenuArea />
        </div>
        <div className={styles.mainDisplayArea}>{props.children}</div>
      </div>
    </>
  )
}

export { Workarea }

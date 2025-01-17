import type { VFC } from 'react'
import styles from '@/components/Topbar/TopbarMenuItemButton.module.css'

type TopbarMenuItemButtonProps = {
  id: string
  children?: React.ReactNode
  title: string
  action: () => void
}

const TopbarMenuItemButton: VFC<TopbarMenuItemButtonProps> = ({ id, title, action }: TopbarMenuItemButtonProps) => {
  return (
    <>
      <div id={id} title={title} className={styles.menuTopItemButton} onClick={action}>
        <div className={styles.menuTopItemTab}></div>
        <div className={styles.menuTopItemTextArea}>
          <div className={styles.menuTopItemText}>{title}</div>
        </div>
      </div>
    </>
  )
}

export { TopbarMenuItemButton }

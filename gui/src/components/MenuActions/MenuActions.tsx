import classNames from 'classnames'
import styles from '@/components/MenuActions/MenuActions.module.css'
import { Cluster } from '@/types/cluster'
import { useEffect } from 'react'
import { MODAL_TYPES, useModalContext } from '@/contexts/ModalContext'

type MenuActionsPosition = {
  top: number
  left: number
}

type MenuActionsProps = {
  isOpen: boolean
  position: MenuActionsPosition
  clusters: Cluster[]
}

const MenuActions = ({ isOpen, position, clusters }: MenuActionsProps) => {
  const { showModal } = useModalContext()

  const clusterStateToggle = (clusters: Cluster[]) => {
    if (!clusters || clusters.length === 0) return

    console.log(clusters.map(cluster => cluster.name))

    if (clusters[0].status === 'Paused') {
      return {
        label: 'Resume UnDistro',
        class: styles.actionsMenuResumeClusterIcon,
        action: () => showModal(MODAL_TYPES.RESUME_CLUSTER, { cluster: clusters[0] })
      }
    } else if (clusters[0].status === 'Ready') {
      return {
        label: 'Pause UnDistro',
        class: styles.actionsMenuPauseClusterIcon,
        action: () => showModal(MODAL_TYPES.PAUSE_CLUSTER, { cluster: clusters[0] })
      }
    }
  }

  const actions = [
    clusterStateToggle(clusters),
    {
      label: 'Update K8s',
      class: styles.actionsMenuUpdateClusterIcon,
      action: () => console.log('update k8s')
    },
    {
      label: 'Cluster Settings',
      class: styles.actionsMenuSettingsClusterIcon,
      action: () => console.log('update k8s')
    },
    {
      label: 'Delete Cluster',
      class: styles.actionsMenuDeleteClusterIcon,
      action: () => showModal(MODAL_TYPES.DELETE_CLUSTER, { cluster: clusters[0] })
    }
  ]

  useEffect(() => {
    console.log(clusters)
  }, [clusters])

  const menuActionsStyles = {
    container: classNames(styles.menuActionsContainer, 'dialogWindowShadow'),
    position: { left: position.left, top: position.top }
  }

  return isOpen ? (
    <div style={menuActionsStyles.position} className={menuActionsStyles.container}>
      <div className={styles.actionsMenu}>
        <ol className={styles.actionsMenuList}>
          {actions.map(item => (
            <li key={`action-${item.label}`} onClick={item.action} className={styles.actionsMenuList}>
              <button className={styles.actionMenuItemContainer}>
                <div className={item.class}></div>
                <a className={styles.actionMenuItemText}>{item.label}</a>
              </button>
            </li>
          ))}
        </ol>
      </div>
      <div className={styles.pointerContainerBorder}>
        <div className={styles.pointerBorder}></div>
      </div>
      <div className={styles.pointerContainer}>
        <div className={styles.pointer}></div>
      </div>
    </div>
  ) : null
}

export { MenuActions }

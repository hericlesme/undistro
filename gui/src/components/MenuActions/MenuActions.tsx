import classNames from 'classnames'
import styles from '@/components/MenuActions/MenuActions.module.css'

type MenuActionsPosition = {
  top: number
  left: number
}

type MenuActionsProps = {
  isOpen: boolean
  position: MenuActionsPosition
}

const MenuActions = ({ isOpen, position }: MenuActionsProps) => {
  const actions = [
    {
      label: 'Pause UnDistro',
      class: styles.actionsMenuPauseClusterIcon
    },
    {
      label: 'Update K8s',
      class: styles.actionsMenuUpdateClusterIcon
    },
    {
      label: 'Cluster Settings',
      class: styles.actionsMenuSettingsClusterIcon
    },
    {
      label: 'Delete Cluster',
      class: styles.actionsMenuDeleteClusterIcon
    }
  ]

  const menuActionsStyles = {
    container: classNames(styles.menuActionsContainer, 'dialogWindowShadow'),
    position: { left: position.left, top: position.top }
  }

  return isOpen ? (
    <div style={menuActionsStyles.position} className={menuActionsStyles.container}>
      <div className={styles.actionsMenu}>
        <ol className={styles.actionsMenuList}>
          {actions.map(action => (
            <li key={`action-${action.label}`} className={styles.actionsMenuList}>
              <button className={styles.actionMenuItemContainer}>
                <div className={action.class}></div>
                <a className={styles.actionMenuItemText}>{action.label}</a>
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

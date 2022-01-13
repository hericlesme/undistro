import menuClustersIcon from '@/public/img/menuClustersIcon.svg'
import menuNodePoolsIcon from '@/public/img/menuNodepoolsIcon.svg'
import menuSecurityIcon from '@/public/img/menuSecurityIcon.svg'
import menuLogsIcon from '@/public/img/menuLogsIcon.svg'

import { LeftMenuItemButton } from '@/components/LeftMenu'
import styles from '@/components/LeftMenu/LeftMenuArea.module.css'

const LeftMenuArea = () => {
  const leftMenuItems = [
    {
      id: 'menuClusterButton',
      alt: 'Clusters',
      src: menuClustersIcon,
      path: '/',
      actions: ['Pause', 'Update K8s', 'Settings', 'Delete']
    },
    {
      id: 'menuNodePoolsButton',
      alt: 'Node Pools',
      src: menuNodePoolsIcon,
      path: '/nodepools',
      actions: ['Create', 'Settings', 'Delete']
    },
    {
      id: 'menuSecurityButton',
      alt: 'Security',
      src: menuSecurityIcon,
      path: '/security',
      actions: ['Create Roles', 'Assign Roles', 'Manage Profiles']
    },
    {
      id: 'menuLogsButton',
      alt: 'Logs',
      src: menuLogsIcon,
      path: '/logs',
      actions: []
    }
  ]

  return (
    <>
      <div className={styles.leftNav}>
        {leftMenuItems.map(item => (
          <LeftMenuItemButton id={item.id} path={item.path} key={`menu-${item.id}`} title={item.alt} item={item} />
        ))}
      </div>
    </>
  )
}

export { LeftMenuArea }

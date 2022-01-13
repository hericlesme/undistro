import { useEffect, useState } from 'react'
import Image from 'next/image'
import classNames from 'classnames'

import styles from '@/components/LeftMenu/LeftMenuItemButton.module.css'
import { useRouter } from 'next/router'

type LeftMenuItemProps = {
  id: string
  children?: React.ReactNode
  title: string
  path: string
  item: any
}

const LeftMenuItemButton = ({ id, title, item, path }: LeftMenuItemProps) => {
  const [isOpen, setIsOpen] = useState(false)

  const router = useRouter()

  useEffect(() => {
    if (router.pathname === path) {
      setIsOpen(true)
    }
  }, [])

  const toggleState = () => {
    router.push(path, undefined, { shallow: true })
  }

  const handleClick = () => {
    router.push(path, undefined, { shallow: true })
  }

  const leftMenuStyles = {
    button: classNames(styles.leftMenuButton, {
      [styles.leftMenuButtonActive]: isOpen
    }),
    arrow: classNames(styles.leftMenuButtonArrow, {
      [styles.leftMenuButtonArrowOpen]: isOpen
    }),
    panel: classNames(styles.leftMenuPanelCollapse, {
      [styles.leftMenuPanelClose]: !isOpen
    })
  }
  return (
    <div id={id} title={title} className={styles.leftMenuButtonContainer}>
      <button onClick={handleClick} className={leftMenuStyles.button}>
        <div className={styles.leftMenuButton}>
          <div className={styles.leftMenuButtonIcon}>
            <Image src={item.src} alt={item.alt} />
          </div>
          <div className={styles.leftMenuButtonText}>
            <a className={'upperCase'}>{title}</a>
          </div>
          {item.actions.length > 0 && <div className={leftMenuStyles.arrow} />}
        </div>
      </button>
      <div className={leftMenuStyles.panel}>
        <ol className={styles.leftMenuPanelList}>
          {item.actions.map((action: string) => (
            <li key={`action-${id}-${action}`} className={styles.leftMenuPanelListItem}>
              {action}
            </li>
          ))}
        </ol>
      </div>
    </div>
  )
}

export { LeftMenuItemButton }

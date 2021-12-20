import type { VFC } from 'react'
import { signOut } from 'next-auth/react'
import Image from 'next/image'
import classNames from 'classnames'

import topBarLogo from '@/public/img/logo-topbar.svg'
import { isIdentityEnabled } from '@/helpers/identity'
import { Navbar, TopbarMenuItemButton } from '@/components/Topbar'

import styles from '@/components/Topbar/Topbar.module.css'

const Topbar: VFC = () => {
  const topbarStyles = {
    container: classNames(styles.topBarContainer, 'responsiveWidth'),
    menu: classNames(styles.topBarMenuArea, 'responsiveWidth')
  }

  const topBarMenuItems = [
    { title: 'create', id: styles.menuCreateButton },
    { title: 'modify', id: styles.menuModifyButton },
    { title: 'manage', id: styles.menuManageButton },
    { title: 'preferences', id: styles.menuPreferencesButton },
    { title: 'about', id: styles.menuAboutButton }
  ]

  const handleLogout = () => {
    signOut({ redirect: true, callbackUrl: '/login' })
  }

  return (
    <header>
      <div className={topbarStyles.container}>
        <div className={styles.topBarLogoArea}>
          <div className={styles.topLogo}>
            <Image src={topBarLogo} alt="UnDistro Logo" />
          </div>
        </div>
        <div className={styles.topBarDividerArea}>
          <div className={styles.topBarDivider}></div>
        </div>
        <div className={topbarStyles.menu}>
          {topBarMenuItems.map((item, index) => (
            <TopbarMenuItemButton key={index} title={item.title} id={item.id} />
          ))}
        </div>
        {isIdentityEnabled() && (
          <div className={styles.logoutArea}>
            <div onClick={handleLogout} className={styles.logoutMenu}>
              <a className={styles.logoutText}>logout</a>
            </div>
          </div>
        )}
      </div>
      <Navbar />
    </header>
  )
}

export { Topbar }

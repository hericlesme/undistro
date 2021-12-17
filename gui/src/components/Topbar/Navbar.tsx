import type { VFC } from 'react'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import Link from 'next/link'
import classNames from 'classnames'

import { useClusters } from '@/contexts/ClusterContext'

import styles from '@/components/Topbar/Navbar.module.css'
export interface Breadcrumb {
  breadcrumb: string
  href: string
}

const Navbar: VFC = () => {
  const router = useRouter()
  const { clusters } = useClusters()
  const [breadcrumbs, setBreadcrumbs] = useState<Array<Breadcrumb> | null>(null)
  const [clustersName, setClustersName] = useState<string[]>(clusters)

  useEffect(() => {
    if (clusters) {
      setClustersName(clusters)
    }
    if (router) {
      const linkPath = router.pathname.split('/')
      linkPath.shift()
      const pathArray = linkPath.map((path, i) => {
        return {
          breadcrumb: path,
          href: '/' + linkPath.slice(0, i + 1).join('/')
        }
      })
      setBreadcrumbs(pathArray)
    }
  }, [router, clusters])

  if (!breadcrumbs) {
    return null
  }

  const navbarStyles = {
    container: classNames(styles.navbarContainer, 'responsiveWidth'),
    breadcrumb: classNames(styles.breadCrumb, 'upperCase'),
    breadcrumbArea: classNames(styles.navbarBreadCrumbArea, 'responsiveWidth')
  }

  let selectedMessage = `multiple clusters selected`
  if (clustersName?.length === 1) {
    selectedMessage = clustersName[0]
  } else if (clustersName?.length === 0 || clustersName == undefined) {
    selectedMessage = 'Select a cluster to begin'
  }

  const renderBreadcrumb = (breadcrumb: Breadcrumb, index: number) => {
    if (index == 0 && breadcrumb.breadcrumb == '') {
      return
    }
    if (index === breadcrumbs.length - 1) {
      return (
        <li key={index + 3}>
          <a>{breadcrumb.breadcrumb}</a>
        </li>
      )
    } else {
      return (
        <li>
          <Link href={breadcrumb.href}>
            <a>{breadcrumb.breadcrumb}</a>
          </Link>
        </li>
      )
    }
  }

  return (
    <>
      <div className={navbarStyles.container}>
        <Link href="/">
          <a className={styles.navbarHomeButtonArea}>
            <div className={styles.navbarHomeIconArea}></div>
          </a>
        </Link>
        <div className={navbarStyles.breadcrumbArea}>
          <ol className={navbarStyles.breadcrumb}>
            <li key="1" className={styles.breadCrumb}>
              <Link href="/">
                <a className={styles.breadCrumb}>clusters</a>
              </Link>
            </li>
            <li key="2">
              <a className={styles.breadCrumbSelObject}>{selectedMessage}</a>
            </li>
            {breadcrumbs.map((breadcrumb, index) => {
              renderBreadcrumb(breadcrumb, index)
            })}
          </ol>
        </div>
        <div className={styles.navbarSearchArea}>
          <input id="searchClear" className={styles.navbarSearchBox} type="search"></input>
          <div className={styles.navbarSearchBoxIcon}></div>
        </div>
      </div>
    </>
  )
}

export { Navbar }

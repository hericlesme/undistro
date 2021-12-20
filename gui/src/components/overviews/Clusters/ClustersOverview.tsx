import { createRef, useEffect, useState, useCallback, VFC } from 'react'
import { useRouter } from 'next/router'
import { useResizeDetector } from 'react-resize-detector'

import { MenuActions } from '@/components/MenuActions/MenuActions'
import { ContentNotFound } from '@/components/ContentNotFound'
import { ClustersOverviewRow, ClustersOverviewEmptyRow, ClustersOverviewFooter } from '@/components/overviews/Clusters'

import { useClusters } from '@/contexts/ClusterContext'
import { paginate } from '@/helpers/pagination'
import { Cluster } from '@/lib/cluster'
import { useFetch } from '@/hooks/query'

import styles from '@/components/overviews/Clusters/ClustersOverview.module.css'
import { ClusterCreation } from './Creation/ClusterCreation'

type ClusterOverviewProps = {
  page: string
}

const ClustersOverview: VFC<ClusterOverviewProps> = ({ page }: ClusterOverviewProps) => {
  const router = useRouter()

  const { clusters: selectedClusters, setClusters: setSelectedClusters } = useClusters()

  const rowHeight = 36
  const tableContainerRef = createRef<HTMLDivElement>()
  const tableRef = createRef<HTMLTableElement>() || undefined
  const [clustersList, setClustersList] = useState<Cluster[]>([])
  const [allChecked, setAllChecked] = useState<boolean>(false)
  const { height } = useResizeDetector<HTMLDivElement>({
    targetRef: tableContainerRef
  })
  const [qtyPages, setQtyPages] = useState<number>(1)
  const [isValidPage, setValidPage] = useState<boolean>(true)
  const [pageSize, setPageSize] = useState<number>(0)
  const [initialContainerSize, setInitialContainerSize] = useState<number>(0)

  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false)
  const [menuPosition, setMenuPosition] = useState({ left: 0, top: 0 })

  const [isCreationOpen, setIsCreationOpen] = useState<boolean>(true)

  const columns = ['provider', 'flavor', 'k8s version', 'cluster group', 'machines', 'age', 'status']

  const { data: clusters, isLoading } = useFetch<Cluster[]>('/api/clusters')

  const changeCheckbox = (checked: boolean) => {
    if (checked) {
      let cls: string[] = []
      clustersList?.forEach(c => {
        cls.push(c.name)
      })
      setSelectedClusters(cls)
    } else {
      setSelectedClusters([])
      setAllChecked(false)
    }
  }

  const handleClick = event => {
    const { target } = event
    const targetRect = target.getBoundingClientRect()

    const tableContainer = document.getElementById('tableContainer')
    const tableContainerRect = tableContainer.getBoundingClientRect()
    const pointerOffset = 4

    setIsMenuOpen(true)
    setMenuPosition({
      left: targetRect.left - tableContainerRect.left + pointerOffset,
      top: targetRect.bottom - tableContainerRect.top + pointerOffset
    })
  }

  const handleUserClick = useCallback(event => {
    event.stopPropagation()
    if (event.target.className.includes('actions')) {
      handleClick(event)
    } else {
      setIsMenuOpen(false)
    }
  }, [])

  useEffect(() => {
    window.addEventListener('click', handleUserClick)
    return () => {
      window.removeEventListener('click', handleUserClick)
    }
  }, [handleUserClick])

  let pageNumber = parseInt(page)

  const pagesCalc = useCallback(() => {
    if (isLoading) {
      setClustersList([])
      return
    }

    if (qtyPages && pageNumber > qtyPages) {
      setValidPage(false)
    } else {
      setValidPage(true)
    }

    const pageFooterHeight = 44
    let pageQtyItems = Math.trunc((height! - pageFooterHeight!) / rowHeight)
    pageQtyItems = pageQtyItems - 1 // ignore table header
    setPageSize(pageQtyItems)

    if (pageQtyItems > 0) {
      let qtyPages = Math.ceil(clusters?.length! / pageQtyItems)
      setQtyPages(qtyPages)
    }
    let pageLists = paginate(clusters, pageQtyItems)
    if (pageLists.length >= pageNumber) {
      let items = pageLists[pageNumber - 1]
      setClustersList(items)
    }
  }, [height, pageNumber, qtyPages, clusters, isLoading])

  useEffect(() => {
    if (height) {
      if (initialContainerSize == 0) {
        setInitialContainerSize(height!)
        return
      }
      if (initialContainerSize > 0) {
        pagesCalc()
      }
    }
  }, [height, initialContainerSize, isValidPage, qtyPages, router, pagesCalc])

  useEffect(() => {
    setSelectedClusters([])
  }, [pageNumber, setSelectedClusters])

  useEffect(() => {
    let displayedClusters = clustersList.filter(e => e.name != '')
    setAllChecked(selectedClusters?.length == displayedClusters?.length)
  }, [clustersList, selectedClusters])

  const renderClusters = () => {
    let cls = []
    for (let i = 0; i < clustersList.length + (pageSize - clustersList.length); i++) {
      if (clustersList[i] === undefined) {
        cls.push(<ClustersOverviewEmptyRow key={i} />)
      } else {
        cls.push(<ClustersOverviewRow key={i} cluster={clustersList[i]} disabled={false} />)
      }
    }
    return cls
  }

  const renderTable = () => (
    <>
      <table ref={tableRef} id="table" className={styles.clustersOverviewTable}>
        <thead>
          <tr>
            <th>
              <div className={styles.tableCheckboxAllIconContainer}>
                <label className={styles.tableCheckboxControlAll}>
                  <input
                    className={styles.tableCheckboxAll}
                    onChange={e => changeCheckbox(e.target.checked)}
                    type="checkbox"
                    name="checkbox"
                    checked={allChecked}
                  />
                </label>
              </div>
            </th>
            <th>
              <div onClick={handleClick} className={styles.tableIconCol}>
                <div className={styles.actionsIconAllDisabled}></div>
              </div>
            </th>
            <th className={styles.responsiveTh}>
              <div className={styles.tableHeaderTitle}>clusters</div>
            </th>
            {columns.map((column, i) => (
              <th key={i}>
                <div className={styles.tableHeaderTitle}>{column}</div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>{renderClusters()}</tbody>
      </table>
      <ClusterCreation isOpen={isCreationOpen} />
      <MenuActions isOpen={isMenuOpen} position={menuPosition} />
      <ClustersOverviewFooter
        total={clusters?.length || 0}
        currentPage={pageNumber}
        qtyPages={qtyPages}
        pageSize={pageSize}
      />
    </>
  )

  const renderPage = () => {
    return isValidPage ? renderTable() : <ContentNotFound />
  }

  return (
    <>
      <div className={styles.clustersOverviewContainer}>
        <div id="tableContainer" className={styles.clustersOverviewTableContainer} ref={tableContainerRef}>
          {renderPage()}
        </div>
      </div>
    </>
  )
}

export { ClustersOverview }

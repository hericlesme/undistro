import type { Cluster } from '@/types/cluster'

import { createRef, useEffect, useState, useCallback, VFC } from 'react'
import { useRouter } from 'next/router'
import { useResizeDetector } from 'react-resize-detector'

import { MenuActions } from '@/components/MenuActions/MenuActions'
import { ContentNotFound } from '@/components/ContentNotFound'
import {
  NodepoolsOverviewRow,
  ClustersOverviewEmptyRow,
  ClustersOverviewFooter
} from '@/components/overviews/MachinePools'
import { useClusters } from '@/contexts/ClusterContext'
import { paginate } from '@/helpers/pagination'
import { useFetch } from '@/hooks/query'

import styles from '@/components/overviews/Clusters/ClustersOverview.module.css'

type NodePoolsOverviewProps = {
  page: string
}

const NodePoolsOverview: VFC<NodePoolsOverviewProps> = ({ page }: NodePoolsOverviewProps) => {
  const router = useRouter()

  const { clusters: selectedClusters } = useClusters()

  const rowHeight = 36
  const tableContainerRef = createRef<HTMLDivElement>()
  const tableRef = createRef<HTMLTableElement>() || undefined
  const [nodepoolsList, setNodepoolsList] = useState<Cluster[]>([])
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

  const columns = ['type', 'replicas', 'k8s version', 'labels', 'taints', 'age', 'status']

  const { data: clusters, isLoading } = useFetch<Cluster[]>('/api/clusters')

  const nodepools =
    clusters &&
    clusters
      .filter(c => selectedClusters.includes(c.name))
      .flatMap(cluster => {
        const controlPlane = {
          ...cluster.controlPlane,
          age: cluster.age,
          labels: 0,
          name: cluster.name,
          status: cluster.status,
          taints: 0,
          type: 'Control Plane',
          replicas: cluster.machines,
          version: cluster.k8sVersion
        }

        const workers = cluster.workers.map((worker: any, i: number) => ({
          ...worker,
          age: cluster.age,
          labels: Object.keys(worker.labels || {}).length,
          name: `${cluster.name}-mp-${i}`,
          status: cluster.status,
          taints: (worker.taints || []).length,
          type: worker.infraNode ? 'InfraNode' : 'Worker',
          replicas: cluster.machines,
          version: cluster.k8sVersion
        }))

        return [controlPlane, ...workers]
      })

  const changeCheckbox = (checked: boolean) => {
    if (checked) {
      let cls: string[] = []
      nodepoolsList?.forEach(c => {
        cls.push(c.name)
      })
    } else {
      setAllChecked(false)
    }
  }

  const getClustersByName = (clusterNames: string[]) => {
    return clusters?.filter(c => clusterNames.includes(c.name))
  }

  const toggleClusterSelection = clusterName => event => {
    event.stopPropagation()
    if (event.ctrlKey) {
      let cls: string[] = [...selectedClusters]
      if (cls.includes(clusterName)) {
        cls = cls.filter(c => c !== clusterName)
      } else {
        cls.push(clusterName)
      }
    }

    if (event.target.className.includes('actions') && selectedClusters.length > 0) {
      handleClick(event)
    } else {
      setIsMenuOpen(false)
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
      setNodepoolsList([])
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
    let pageLists = paginate(nodepools, pageQtyItems)
    if (pageLists.length >= pageNumber) {
      let items = pageLists[pageNumber - 1]
      setNodepoolsList(items)
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

  const renderClusters = () => {
    let cls = []
    for (let i = 0; i < nodepoolsList.length + (pageSize - nodepoolsList.length); i++) {
      if (nodepoolsList[i] === undefined) {
        cls.push(<ClustersOverviewEmptyRow key={i} />)
      } else {
        cls.push(
          <NodepoolsOverviewRow
            onClick={toggleClusterSelection(nodepoolsList[i].name)}
            key={i}
            nodepool={nodepoolsList[i]}
            disabled={false}
          />
        )
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
              <div className={styles.tableIconCol}>
                <div className={styles.actionsIconAllDisabled}></div>
              </div>
            </th>
            <th className={styles.responsiveTh}>
              <div className={styles.tableHeaderTitle}>Name</div>
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
      <MenuActions isOpen={isMenuOpen} position={menuPosition} clusters={getClustersByName(selectedClusters)} />
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

export { NodePoolsOverview as ClustersOverview }

import type { VFC } from 'react'

import { forwardRef } from 'react'
import { useRouter } from 'next/router'

import Pagination from '@/components/Pagination/Pagination'

import styles from '@/components/overviews/Clusters/ClustersOverviewNavFooter.module.css'

type ClustersOverviewFooterProps = {
  total: number
  currentPage: number
  pageSize: number
  qtyPages: number
  refer?: React.ForwardedRef<HTMLDivElement>
}

const ClustersOverviewFooter: VFC<ClustersOverviewFooterProps> = forwardRef<
  HTMLDivElement,
  ClustersOverviewFooterProps
>((props, ref) => (
  <ClustersOverviewNavFooter
    total={props.total || 0}
    currentPage={props.currentPage}
    qtyPages={props.qtyPages}
    pageSize={props.pageSize}
    refer={ref}
  />
))

ClustersOverviewFooter.displayName = 'ClusterOverviewFooter'

const ClustersOverviewNavFooter: VFC<ClustersOverviewFooterProps> = (props: ClustersOverviewFooterProps) => {
  const router = useRouter()

  const pressEnter = (e: React.KeyboardEvent<HTMLInputElement>) => {
    let el = e.target as HTMLInputElement
    router.query.page = el.value
    if (e.key == 'Enter' && parseInt(el.value) != props.currentPage) {
      if (parseInt(el.value) <= props.qtyPages) {
        router.push({
          query: { page: el.value }
        })
      }
    }
  }

  const onBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    let el = e.target as HTMLInputElement
    if (parseInt(el.value) <= props.qtyPages && parseInt(el.value) != props.currentPage) {
      router.push({
        query: { page: el.value }
      })
    }
  }

  return (
    <div ref={props.refer} id="pageFooter" className={styles.tableFooterContainer}>
      <div className={styles.tableFooter}>
        <div className={styles.navFooterResults}>
          <a className={styles.navFooterResultsText}>
            {props.total} {props.total > 1 ? 'Results' : 'Result'}
          </a>
        </div>
        {props.qtyPages > 1 && (
          <>
            <div className={styles.navFooterJumpToPage}>
              <a className={styles.navFooterJumpToPageText}>Jump to page</a>
            </div>

            <div className={styles.paginationSearchArea}>
              <input
                onKeyPress={e => pressEnter(e)}
                onBlur={e => onBlur(e)}
                className={styles.paginationSearchBox}
                placeholder={props.currentPage.toString()}
                type="text"
              ></input>
            </div>
            <Pagination
              currentPage={props.currentPage}
              totalCount={props.total}
              pageSize={props.pageSize}
              onPageChange={page => console.log(page)}
            />
          </>
        )}
      </div>
    </div>
  )
}

export { ClustersOverviewFooter }

import type { VFC } from 'react'
import styles from '@/components/overviews/NodePools/NodePoolsOverviewNavFooter.module.css'

const NodePoolsOverviewNavFooter: VFC = () => (
  <div className={styles.tableFooterContainer}>
    <div className={styles.tableFooter}>
      <div className={styles.navFooterResults}>
        <a className={styles.navFooterResultsText}>10 Results</a>
      </div>

      <div className={styles.navFooterJumpToPage}>
        <a className={styles.navFooterJumpToPageText}>Jump to page</a>
      </div>

      <div className={styles.paginationSearchArea}>
        <input className={styles.paginationSearchBox} placeholder="5" type="text"></input>
      </div>

      <div className={styles.paginationNavContainer}>
        <a href="#">
          <div className={styles.paginationNavArrowLeft}></div>
        </a>
        <a className={styles.paginationNavPagesText}>1</a>
        <a className={styles.paginationNavPagesInterval}>...</a>
        <a className={styles.paginationNavPagesText}>4</a>
        <a className={styles.paginationNavCurrentPage}>5</a>
        <a className={styles.paginationNavPagesText}>6</a>
        <a className={styles.paginationNavPagesInterval}>...</a>
        <a className={styles.paginationNavPagesText}>10</a>
        <a href="#">
          <div className={styles.paginationNavArrowRight}></div>
        </a>
      </div>
    </div>
  </div>
)

export { NodePoolsOverviewNavFooter }

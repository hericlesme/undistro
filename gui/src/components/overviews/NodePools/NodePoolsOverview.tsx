import type { VFC } from 'react'

import { NodePoolsOverviewInfranodeRow, NodePoolsOverviewNavFooter } from '@/components/overviews/NodePools'

import styles from '@/components/overviews/NodePools/NodepoolsOverview.module.css'

const NodePoolsOverview: VFC = () => (
  <div className={styles.clustersOverviewContainer}>
    <div className={styles.clustersOverviewTableContainer}>
      <table className={styles.clustersOverviewTable}>
        <thead>
          <tr>
            <th>
              <div className={styles.tableCheckboxAllIconContainer}>
                <label className={styles.tableCheckboxControlAll}>
                  <input className={styles.tableCheckboxAll} type="checkbox" name="checkbox" />
                </label>
              </div>
            </th>
            <th>
              <div className={styles.tableIconCol}>
                <div className={styles.actionsIconAllDisabled}></div>
              </div>
            </th>
            <th className={styles.responsiveTh}>
              <div className={styles.tableHeaderTitle}>nodepool</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>type</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>replicas</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>k8s version</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>labels</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>taints</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>age</div>
            </th>
            <th>
              <div className={styles.tableHeaderTitle}>status</div>
            </th>
          </tr>
        </thead>
        <tbody>
          <NodePoolsOverviewInfranodeRow />
        </tbody>
      </table>
    </div>
    <NodePoolsOverviewNavFooter />
  </div>
)

export { NodePoolsOverview }

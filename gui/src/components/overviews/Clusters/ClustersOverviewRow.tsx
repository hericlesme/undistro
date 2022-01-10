import type { VFC } from 'react'
import type { Cluster } from '@/types/cluster'

import classNames from 'classnames'

import { useClusters } from '@/contexts/ClusterContext'

import styles from '@/components/overviews/Clusters/ClustersOverviewRow.module.css'

type ClustersOverviewRowProps = {
  cluster: Cluster
  disabled: boolean
  onClick?: (event: any) => void
}

const ClustersOverviewRow: VFC<ClustersOverviewRowProps> = ({
  cluster,
  disabled,
  ...props
}: ClustersOverviewRowProps) => {
  const { clusters, setClusters } = useClusters()

  const tableStyles = {
    uppercase: classNames(styles.tableCellTitle, 'upperCase'),
    centered: {
      default: classNames(styles.tableCellTitle, 'textCentered'),
      uppercase: classNames(styles.tableCellTitle, 'upperCase', 'textCentered'),
      warning: classNames(styles.tableCellTitleWarning, 'textCentered'),
      critical: classNames(styles.tableCellTitleWarning, 'textCentered'),
      success: classNames(styles.tableCellTitleSuccess, 'textCentered')
    }
  }

  let statusClass = tableStyles.centered.critical

  if (cluster.status.toLowerCase() == 'ready') {
    statusClass = tableStyles.centered.success
  } else if (
    cluster.status.toLowerCase() == 'provisioning' ||
    cluster.status.toLowerCase() == 'paused' ||
    cluster.status.toLowerCase() == 'deleting'
  ) {
    statusClass = tableStyles.centered.warning
  }

  const changeCheckbox = (name: string, checked: boolean) => {
    let cls = new Set<string>(clusters)
    if (checked) {
      cls.add(name)
    } else if (cls.has(name)) {
      cls.delete(name)
    }
    setClusters(Array.from(cls.values()))
  }

  const handleCheckBoxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    changeCheckbox(cluster.name, event.target.checked)
  }

  let machineStr = ''
  if (!disabled) {
    machineStr = cluster.machines.toString()
  }

  return (
    <tr {...props}>
      <td>
        <div className={styles.tableCheckboxIconContainer}>
          <label className={styles.tableCheckboxControl}>
            <input
              className={styles.tableCheckbox}
              onChange={handleCheckBoxChange}
              type="checkbox"
              name="checkbox"
              checked={clusters.includes(cluster.name)}
              disabled={disabled}
            />
          </label>
        </div>
      </td>
      <td>
        <div className={styles.tableActionsIconContainer}>
          <div className={styles.actionsIcon}></div>
        </div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{cluster.name}</div>
      </td>
      <td>
        <div className={tableStyles.uppercase}>{cluster.provider}</div>
      </td>
      <td>
        <div className={tableStyles.centered.uppercase}>{cluster.flavor}</div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{cluster.k8sVersion}</div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{cluster.clusterGroup}</div>
      </td>
      <td>
        <div className={tableStyles.centered.default}>{machineStr}</div>
      </td>
      <td>
        <div className={tableStyles.centered.default}>{cluster.age}</div>
      </td>
      <td>
        <div className={statusClass}>{cluster.status}</div>
      </td>
    </tr>
  )
}

export { ClustersOverviewRow }

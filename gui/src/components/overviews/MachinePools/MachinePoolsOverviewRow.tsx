import type { VFC } from 'react'
import type { Cluster } from '@/types/cluster'

import classNames from 'classnames'

import { useClusters } from '@/contexts/ClusterContext'

import styles from '@/components/overviews/Clusters/ClustersOverviewRow.module.css'

type NodepoolsOverviewRowProps = {
  nodepool: any
  disabled: boolean
  onClick?: (event: any) => void
}

const NodepoolsOverviewRow: VFC<NodepoolsOverviewRowProps> = ({
  nodepool,
  disabled,
  ...props
}: NodepoolsOverviewRowProps) => {
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

  if (nodepool.status.toLowerCase() == 'ready') {
    statusClass = tableStyles.centered.success
  } else if (
    nodepool.status.toLowerCase() == 'provisioning' ||
    nodepool.status.toLowerCase() == 'paused' ||
    nodepool.status.toLowerCase() == 'deleting'
  ) {
    statusClass = tableStyles.centered.warning
  }

  const handleCheckBoxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    // changeCheckbox(nodepool.name, event.target.checked)
  }

  let machineStr = ''
  if (!disabled) {
    machineStr = nodepool.taints.toString()
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
        <div className={styles.tableCellTitle}>{nodepool.name}</div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{nodepool.type}</div>
      </td>
      <td>
        <div className={tableStyles.centered.uppercase}>{nodepool.replicas}</div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{nodepool.version}</div>
      </td>
      <td>
        <div className={styles.tableCellTitle}>{nodepool.labels}</div>
      </td>
      <td>
        <div className={tableStyles.centered.default}>{machineStr}</div>
      </td>
      <td>
        <div className={tableStyles.centered.default}>{nodepool.age}</div>
      </td>
      <td>
        <div className={statusClass}>{nodepool.status}</div>
      </td>
    </tr>
  )
}

export { NodepoolsOverviewRow }

import type { VFC } from 'react'
import styles from '@/components/overviews/NodePools/NodePoolsOverviewInfranodeRow.module.css'

const NodePoolsOverviewInfranodeRow: VFC = () => {
  let tableCellTitleCentered = [styles.tableCellTitle, 'textCentered'].join(' ')
  let tableCellTitleUpperCaseCentered = [styles.tableCellTitle, 'upperCase', 'textCentered'].join(' ')
  let tableCellTitleWarningCentered = [styles.tableCellTitleWarning, 'textCentered'].join(' ')
  let tableCellTitleSuccessCentered = [styles.tableCellTitleSuccess, 'textCentered'].join(' ')
  let tableCellTitleCriticalCentered = [styles.tableCellTitleCritical, 'textCentered'].join(' ')
  return (
    <>
      <tr>
        <td>
          <div className={styles.tableCheckboxIconContainer}>
            <label className={styles.tableCheckboxControl}>
              <input className={styles.tableCheckbox} type="checkbox" name="checkbox" />
            </label>
          </div>
        </td>
        <td>
          <div className={styles.tableActionsIconContainer}>
            <div className={styles.actionsIcon}></div>
          </div>
        </td>
        <td>
          <div className={styles.tableCellTitle}>nodepool-01</div>
        </td>
        <td>
          <div className={styles.tableCellTitle}>infranode</div>
        </td>
        <td>
          <div className={tableCellTitleCentered}>3</div>
        </td>
        <td>
          <div className={styles.tableCellTitle}>v.1.17.7</div>
        </td>
        <td>
          <div className={tableCellTitleCentered}>1</div>
        </td>
        <td>
          <div className={tableCellTitleCentered}>2</div>
        </td>
        <td>
          <div className={tableCellTitleCentered}>2d21h</div>
        </td>
        <td>
          <div className={tableCellTitleWarningCentered}>paused</div>
        </td>
      </tr>
    </>
  )
}

export { NodePoolsOverviewInfranodeRow }

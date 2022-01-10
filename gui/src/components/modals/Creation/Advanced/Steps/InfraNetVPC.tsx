import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import { useFetch } from '@/hooks/query'

import { Provider } from '@/types/cluster'
import { useWatch } from 'react-hook-form'
import { TextInput, Select } from '@/components/forms'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import classNames from 'classnames'

const InfraNetVPC: VFC<FormActions> = ({ register }: FormActions) => {
  return (
    <>
      <div className={styles.infraNetworkInputRowSmall}>
        <div className={classNames(styles.inputBlock, styles.inputFit)}>
          <label className={styles.createClusterLabel} htmlFor="infraNetworkID">
            ID
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="infraNetworkID"
            name="infraNetworkID"
          >
            <option value="" disabled selected hidden>
              default
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
      <div className={classNames(styles.infraNetworkInputRowSmall, styles.rowMargin)}>
        <div className={classNames(styles.inputBlock, styles.inputSmall)}>
          <label className={styles.createClusterLabel} htmlFor="infraNetworkZone">
            zone
          </label>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="zone"
            id="infraNetworkZone"
            name="infraNetworkZone"
          />
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={classNames(styles.inputBlock, styles.inputFit)}>
          <label className={styles.createClusterLabel} htmlFor="infraNetworkZoneCidrBlock">
            CIDR block
          </label>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="optional"
            id="infraNetworkZoneCidrBlock"
            name="infraNetworkZoneCidrBlock"
          />
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
      <div className={classNames(styles.modalSubnetContainer, styles.bordered)}>
        <div className={styles.workersTitleContainer}>
          <a className={styles.modalCreateClusterTitle}>subnet</a>
        </div>
        <div className={styles.modalSubnetBlock}>
          <div className={styles.infraNodeBlock}>
            <div className={classNames(styles.switchContainer, styles.justifyLeft)}>
              <a className={styles.createClusterLabel}>is public</a>
              <label className={styles.switch} htmlFor="infraNetworkZoneIsPublic">
                <input type="checkbox" id="infraNetworkZoneIsPublic" name="infraNetworkZoneIsPublic" />
                <span className={classNames(styles.slider, styles.round)}></span>
              </label>
            </div>
          </div>
          <div className={styles.subnetInputRow}>
            <div className={classNames(styles.inputBlock, styles.inputFit)}>
              <label className={styles.createClusterLabel} htmlFor="infraNetworkSubnetID">
                ID
              </label>
              <select
                className={classNames(styles.createClusterTextSelect, styles.input100)}
                id="infraNetworkSubnetID"
                name="infraNetworkSubnetID"
              >
                <option value="" disabled selected hidden>
                  default
                </option>
                <option value="option1">option1</option>
                <option value="option2">option2</option>
                <option value="option3">option3</option>
              </select>
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
            <div className={classNames(styles.inputBlock, styles.inputMedium)}>
              <label className={styles.createClusterLabel} htmlFor="infraNetworkSubnetZone">
                zone
              </label>
              <input
                className={classNames(styles.createClusterTextInput, styles.input100)}
                placeholder="zone"
                id="infraNetworkSubnetZone"
                name="infraNetworkSubnetZone"
              />
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
            <div className={classNames(styles.inputBlock, styles.inputMedium)}>
              <label className={styles.createClusterLabel} htmlFor="infraNetworkSubnetCidrBlock">
                CIDR block
              </label>
              <input
                className={classNames(styles.createClusterTextInput, styles.input100)}
                placeholder="optional"
                id="infraNetworkSubnetCidrBlock"
                name="infraNetworkSubnetCidrBlock"
              />
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
          </div>
          <div className={classNames(styles.addButtonBlock, styles.justifyRight)}>
            <button className={styles.solidMdButtonDefault}>
              <a>add</a>
            </button>
          </div>

          <div className={classNames(styles.modalTableContainer, styles.bordered)}>
            <table className={styles.modalWorkersTable} id="wizardWorkersTable">
              <tbody>
                <tr>
                  <td>
                    <div className={styles.modalTableRow}>
                      <div className={styles.modalTableItem}>
                        <a>subnet-0</a>
                      </div>
                      <div className={styles.modalDeleteTableItemContainer}>
                        <button className={styles.deleteTableItem}></button>
                      </div>
                    </div>
                  </td>
                </tr>
                <tr>
                  <td>
                    <div className={styles.modalTableRow}>
                      <div className={styles.modalTableItem}>
                        <a>subnet-1</a>
                      </div>
                      <div className={styles.modalDeleteTableItemContainer}>
                        <button className={styles.deleteTableItem}></button>
                      </div>
                    </div>
                  </td>
                </tr>
                <tr>
                  <td>
                    <div className={styles.modalTableRow}>
                      <div className={styles.modalTableItem}>
                        <a>subnet-2</a>
                      </div>
                      <div className={styles.modalDeleteTableItemContainer}>
                        <button className={styles.deleteTableItem}></button>
                      </div>
                    </div>
                  </td>
                </tr>
                <tr>
                  <td>
                    <div className={styles.modalTableRow}>
                      <div className={styles.modalTableItem}>
                        <a>subnet-3</a>
                      </div>
                      <div className={styles.modalDeleteTableItemContainer}>
                        <button className={styles.deleteTableItem}></button>
                      </div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  )
}

export { InfraNetVPC }

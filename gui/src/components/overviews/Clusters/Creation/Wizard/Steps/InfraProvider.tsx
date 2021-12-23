import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'

const InfraProvider: VFC<FormActions> = ({ register }: FormActions) => {
  return (
    <>
      <div className={styles.inputRow}>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="infraProviderID">
            ID
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="infraProviderID"
            name="infraProviderID"
            {...register('infraProviderID', { required: true })}
          >
            <option value="" disabled selected hidden>
              ID
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="infraProviderCIDR">
            CIDR block
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="infraProviderCIDR"
            name="infraProviderCIDR"
            {...register('infraProviderCIDR', { required: true })}
          >
            <option value="" disabled selected hidden>
              CIDR block
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
      <div className={styles.inputRow}>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="infraProviderFlavor">
            flavor
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="infraProviderFlavor"
            name="infraProviderFlavor"
            {...register('infraProviderFlavor', { required: true })}
          >
            <option value="" disabled selected hidden>
              select flavor
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="infraProviderK8sVersion">
            kubernetes version
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="infraProviderK8sVersion"
            name="infraProviderK8sVersion"
            {...register('infraProviderK8sVersion', { required: true })}
          >
            <option value="" disabled selected hidden>
              select K8s
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
      <div className={styles.inputBlockTextArea}>
        <label className={styles.createClusterLabel} htmlFor="infraProviderSshKey">
          sshKey
        </label>
        <br />
        <textarea
          className={classNames(styles.createClusterTextInput, styles.input100, styles.textAreaInput)}
          placeholder="ssh Key"
          id="infraProviderSshKey"
          name="infraProviderSshKey"
          {...register('infraProviderSshKey', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>

      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="infraProviderUploadFile">
          upload config file
        </label>
        <div className={styles.inputUploadFileBlock}>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="clouds.yaml"
            type="text"
            id="infraProviderUploadFile"
            name="infraProviderUploadFile"
            {...register('infraProviderUploadFile', { required: true })}
          />
          <div className={styles.uploadFileButtonContainer}>
            <button className={styles.uploadFileIcon}></button>
          </div>
        </div>
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
    </>
  )
}

export { InfraProvider }

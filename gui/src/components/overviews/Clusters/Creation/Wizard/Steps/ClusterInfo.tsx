import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'

const ClusterInfo: VFC<FormActions> = ({ register }: FormActions) => {
  return (
    <>
      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="clusterName">
          Cluster name
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="choose a cool name for this cluster"
          type="text"
          id="clusterName"
          name="clusterName"
          {...register('clusterName', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="clusterNamespace">
          Namespace
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="namespace"
          type="text"
          id="clusterNamespace"
          name="clusterNamespace"
          {...register('clusterNamespace', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
      <div className={styles.inputRow}>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="clusterProvider">
            Provider
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="clusterProvider"
            name="clusterProvider"
            {...register('clusterProvider', { required: true })}
          >
            <option value="" disabled selected hidden>
              Select provider
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="clusterDefaultRegion">
            Default region
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="clusterDefaultRegion"
            name="clusterDefaultRegion"
            {...register('clusterDefaultRegion', { required: true })}
          >
            <option value="" disabled selected hidden>
              Select region
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
    </>
  )
}

export { ClusterInfo }

import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'

const AddOns: VFC<FormActions> = ({ register }: FormActions) => {
  return (
    <>
      <div className={styles.addOnBlock}>
        <div className={classNames(styles.switchContainer, styles.justifyRight)}>
          <a className={styles.createClusterLabel}>default policies</a>
          <label className={styles.switch} htmlFor="addOnDefaultPolicies">
            <input
              type="checkbox"
              id="addOnDefaultPolicies"
              name="addOnDefaultPolicies"
              {...register('addOnDefaultPolicies')}
            />
            <span className={classNames(styles.slider, styles.round)}></span>
          </label>
        </div>
        <div className={styles.addOnDescriptionTitlesContainer}>
          <a className={styles.addOnDescriptionTitles}>applying default policies will enable cluster best</a>
          <a className={styles.addOnDescriptionTitles}>practices and security policies.</a>
        </div>
      </div>
      <div className={styles.addOnBlock}>
        <div className={classNames(styles.switchContainer, styles.justifyRight)}>
          <a className={styles.createClusterLabel}>observer</a>
          <label className={styles.switch} htmlFor="addOnObserver">
            <input type="checkbox" id="addOnObserver" name="addOnObserver" {...register('addOnObserver')} />
            <span className={classNames(styles.slider, styles.round)}></span>
          </label>
        </div>
        <div className={styles.addOnDescriptionTitlesContainer}>
          <a className={styles.addOnDescriptionTitles}>this will install and configure logging and</a>
          <a className={styles.addOnDescriptionTitles}>monitoring tools.</a>
        </div>
      </div>
      <div className={styles.addOnBlock}>
        <div className={classNames(styles.switchContainer, styles.justifyRight)}>
          <a className={styles.createClusterLabel}>identity</a>
          <label className={styles.switch} htmlFor="addOnIdentity">
            <input type="checkbox" id="addOnIdentity" name="addOnIdentity" {...register('addOnIdentity')} />
            <span className={classNames(styles.slider, styles.round)}></span>
          </label>
        </div>
        <div className={styles.addOnDescriptionTitlesContainer}>
          <a className={styles.addOnDescriptionTitles}>enable centralized authentication system between</a>
          <a className={styles.addOnDescriptionTitles}>all clusters managed by UnDistro.</a>
        </div>
      </div>
    </>
  )
}

export { AddOns }

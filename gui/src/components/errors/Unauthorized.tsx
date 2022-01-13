import classNames from 'classnames'
import * as React from 'react'
import styles from './Unauthorized.module.css'

const Unauthorized = () => {
  return (
    <>
      <div className={styles.page404ContainerMessage}>
        <div className={styles.page404messageContainer}>
          ​<div className={styles.page404MonitorMessage}></div>​<div className={styles.errorTitle}>access denied</div>
          <div className={classNames(styles.errorDescription, styles.warning)}>
            You are not allowed to access this area.
          </div>
          <div className={styles.errorDescription}>Talk to your administrator to unlock this feature.</div>
        </div>
      </div>
    </>
  )
}
export default Unauthorized

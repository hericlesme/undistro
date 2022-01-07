import { useEffect, useState, VFC } from 'react'

import classNames from 'classnames'
import { useModalContext } from '@/contexts/ModalContext'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import { Logs } from './Logs'

type StatusProps = {
  creationInfo: {
    namespace: string
    cluster: string
    status: string
  }
}

const Progress: VFC<StatusProps> = ({ creationInfo }: StatusProps) => {
  const { hideModal } = useModalContext()

  let { namespace, cluster, status } = creationInfo

  console.log(creationInfo)

  return (
    <>
      <div className={styles.createClusterWizContainer}>
        <div className={styles.modalTitleContainer}>
          <a className={styles.modalCreateClusterTitle}>Progress</a>
        </div>
        <div className={styles.modalContentContainer}>
          <div className={styles.modalContentContainer}>
            <div className={classNames(styles.modalProgressArea, styles.modalInputArea)}>
              <div className={styles.modalForm}>
                <div className={styles.modalProgressArea}>
                  <div className={styles.progressBarContainer}>
                    <div className={styles.progressBar}></div>
                  </div>
                  <Logs namespace={namespace} cluster={cluster} />
                </div>
              </div>
            </div>
            <div className={styles.modalDialogButtonsContainer}>
              <div className={styles.rightButtonContainer}>
                <button onClick={hideModal} className={styles.borderButtonSuccess}>
                  <a>Close</a>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export { Progress }

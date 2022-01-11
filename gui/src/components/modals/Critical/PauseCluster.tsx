import { DialogOverlay, DialogContent } from '@reach/dialog'

import styles from '@/components/modals/Critical/CriticalDialog.module.css'
import { Cluster } from '@/types/cluster'
import { useModalContext } from '@/contexts/ModalContext'
import { useMutate } from '@/hooks/query'

type PauseClusterProps = {
  cluster: Cluster
}

const PauseCluster = ({ cluster }: PauseClusterProps) => {
  const { hideModal } = useModalContext()
  const pauseCluster = useMutate({
    url: `/api/clusters/${cluster.clusterGroup}/${cluster.name}/edit`,
    method: 'patch',
    invalidate: '/api/clusters'
  })

  const handlePause = async () => {
    pauseCluster.mutate({
      spec: {
        paused: true
      }
    })
    hideModal()
  }

  return (
    <DialogOverlay isOpen={true} className={styles.criticalDialogOverlay}>
      <DialogContent aria-label="Pause Cluster" className={styles.criticalDialogStripe}>
        <div className={styles.criticalDialogMessageContainer}>
          <div className={styles.criticalDialogObjectName}>{cluster.name}</div>
          <div className={styles.warningDialogIcon}></div>
          <div className={styles.criticalDialogActionDescription}>pause</div>
          <div className={styles.warningDialogActionAlertMessage}>
            This cluster will remain active, but UnDistro management will be stopped. You can resume it at any time.
          </div>
          <div className={styles.warningDialogBtnContainer}>
            <button onClick={handlePause} className={styles.confirmDialogBtn}>
              <div className={styles.confirmBtnText}>pause</div>
            </button>

            <button onClick={hideModal} className={styles.cancelDialogBtn}>
              <div className={styles.cancelBtnText}>cancel</div>
            </button>
          </div>
        </div>
      </DialogContent>
    </DialogOverlay>
  )
}

export { PauseCluster }

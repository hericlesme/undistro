import { DialogOverlay, DialogContent } from '@reach/dialog'

import styles from '@/components/modals/Critical/CriticalDialog.module.css'
import { Cluster } from '@/types/cluster'
import { useModalContext } from '@/contexts/ModalContext'
import { useMutate } from '@/hooks/query'

type ResumeClusterProps = {
  cluster: Cluster
}

const ResumeCluster = ({ cluster }: ResumeClusterProps) => {
  const { hideModal } = useModalContext()
  const pauseCluster = useMutate({
    url: `/api/clusters/${cluster.clusterGroup}/${cluster.name}/edit`,
    method: 'patch',
    invalidate: '/api/clusters'
  })

  const handlePause = async () => {
    console.log(cluster.clusterGroup, cluster.name)
    pauseCluster.mutate({
      spec: {
        paused: false
      }
    })
    hideModal()
  }

  return (
    <DialogOverlay isOpen={true} className={styles.criticalDialogOverlay}>
      <DialogContent aria-label="Pause Cluster" className={styles.criticalDialogStripe}>
        <div className={styles.criticalDialogMessageContainer}>
          <div className={styles.criticalDialogObjectName}>{cluster.name}</div>
          <div className={styles.successDialogIcon}></div>
          <div className={styles.criticalDialogActionDescription}>resume</div>
          <div className={styles.successDialogActionAlertMessage}>Resume UnDistro cluster management.</div>
          <div className={styles.criticalDialogBtnContainer}>
            <button onClick={handlePause} className={styles.confirmDialogBtn}>
              <div className={styles.confirmBtnText}>resume</div>
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

export { ResumeCluster }

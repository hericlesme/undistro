import { DialogOverlay, DialogContent } from '@reach/dialog'

import styles from '@/components/modals/Critical/CriticalDialog.module.css'
import { Cluster } from '@/types/cluster'
import { useModalContext } from '@/contexts/ModalContext'
import { useState } from 'react'
import { useMutate } from '@/hooks/query'

type DeleteClusterProps = {
  cluster: Cluster
}

const DeleteCluster = ({ cluster }: DeleteClusterProps) => {
  const [deleteConfirm, setDeleteConfirm] = useState('')

  const { hideModal } = useModalContext()
  const deleteCluster = useMutate({
    url: `/api/clusters/${cluster.clusterGroup}/${cluster.name}/delete`,
    method: 'delete',
    invalidate: '/api/clusters'
  })

  const handleDelete = async () => {
    console.log(cluster.clusterGroup, cluster.name)
    deleteCluster.mutate({})
    hideModal()
  }

  return (
    <DialogOverlay isOpen={true} className={styles.criticalDialogOverlay}>
      <DialogContent aria-label="Create Cluster" className={styles.criticalDialogStripe}>
        <div className={styles.criticalDialogMessageContainer}>
          <div className={styles.criticalDialogObjectName}>{cluster.name}</div>
          <div className={styles.criticalDialogIcon}></div>
          <div className={styles.criticalDialogActionDescription}>delete cluster</div>
          <div className={styles.criticalDialogActionAlertMessage}>This action cannot be undone.</div>
          <div className={styles.criticalDialogActionConfirmationContainer}>
            <div className={styles.criticalDialogConfirmationText}>Type “delete” to confirm:</div>
            <input
              id="searchClear"
              onChange={event => setDeleteConfirm(event.target.value)}
              className={styles.confirmationInputBox}
              type="text"
              placeholder="delete"
            ></input>
          </div>
          <div className={styles.criticalDialogBtnContainer}>
            <button
              onClick={handleDelete}
              className={styles.confirmDialogBtn}
              disabled={deleteConfirm.toLowerCase() !== 'delete'}
            >
              <div className={styles.confirmBtnText}>delete</div>
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

export { DeleteCluster }

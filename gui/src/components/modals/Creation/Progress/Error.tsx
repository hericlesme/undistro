import styles from '@/components/modals/Creation/ClusterCreation.module.css'

const CreationError = () => {
  return (
    <div className={styles.progressDialogContainer}>
      <div className={styles.errorMessageContainer}>
        <div className={styles.progressResultTextContainer}>
          <a className={styles.progressResultText}>THE CLUSTER COULD</a>
          <a className={styles.progressResultText}>NOT BE CREATED</a>
        </div>
        <div className={styles.createClusterErrorIcon}></div>
        <div className={styles.progressResultMessageContainer}>
          <a className={styles.progressErrorMessage}>An error occurred while creating the cluster.</a>
          <a className={styles.progressErrorCode}>### - ERROR CODE DESCRIPTION TEXT LINE</a>
        </div>
      </div>
    </div>
  )
}
export { CreationError }

import { AlertDialogDescription } from '@reach/alert-dialog'
import { ReactComponent as DefaultIcon } from '@assets/icons/cluster-actions/default.svg'

type Props = {
  currentVersion: string
}

export const UpgradeClusterAlertContent = ({ currentVersion }: Props) => {
  return (
    <AlertDialogDescription>
      <DefaultIcon />
      <h3 className="action">Upgrade Cluster</h3>
      <p className="warning-message">
        Upgrade to {currentVersion} - This action cannot be undone.
      </p>
    </AlertDialogDescription>
  )
}

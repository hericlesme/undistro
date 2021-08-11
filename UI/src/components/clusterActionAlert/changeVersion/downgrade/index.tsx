import { AlertDialogDescription } from '@reach/alert-dialog'
import { ReactComponent as CriticalIcon } from '@assets/icons/cluster-actions/critical.svg'

type Props = {
  currentVersion: string
}

export const DowngradeClusterAlertContent = ({ currentVersion }: Props) => {
  return (
    <AlertDialogDescription>
      <CriticalIcon />
      <h3 className="action">Downgrade Cluster</h3>
      <p className="danger-message">
        Downgrade to {currentVersion} - This action can cause a cluster failure.
      </p>
    </AlertDialogDescription>
  )
}

import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogLabel
} from '@reach/alert-dialog'
import compareVersions from 'compare-versions'
import { useRef, useState } from 'react'
import { ClusterActionAlertProps } from '../index'
import { ClusterAlertButton } from '../common/button'
import { ClusterAlertSelect } from '../common/select'
import { DowngradeClusterAlertContent } from './downgrade'
import { UpgradeClusterAlertContent } from './upgrade'
import './index.scss'

enum OperationType {
  DOWNGRADE,
  NONE,
  UPGRADE
}

type Props = ClusterActionAlertProps & {
  currentVersion: string
  versions: string[]
  onActionConfirm: (version: string) => void
}

export const ChangeVersionAlert = ({
  currentVersion,
  heading,
  isOpen,
  versions,
  onActionConfirm: handleActionConfirm,
  onDismiss: handleDismiss
}: Props) => {
  const leastDestructiveRef = useRef(null)
  const [selectedVersion, setSelectedVersion] = useState(currentVersion)
  const [operationType, setOperationType] = useState<OperationType>(
    OperationType.NONE
  )

  const defineOperationType = () => {
    const operationType = compareVersions(selectedVersion, currentVersion)

    switch (operationType) {
      case -1:
        setOperationType(OperationType.DOWNGRADE)
        break
      case 0:
        setOperationType(OperationType.NONE)
        break
      case 1:
        setOperationType(OperationType.UPGRADE)
        break
    }
  }

  return (
    <>
      {isOpen && (
        <AlertDialog leastDestructiveRef={leastDestructiveRef}>
          <AlertDialogLabel>{heading}</AlertDialogLabel>

          {operationType === OperationType.DOWNGRADE && (
            <DowngradeClusterAlertContent currentVersion={selectedVersion} />
          )}

          {operationType === OperationType.UPGRADE && (
            <UpgradeClusterAlertContent currentVersion={selectedVersion} />
          )}

          {operationType === OperationType.NONE && (
            <>
              <AlertDialogDescription>
                <p className="change-version-sub-heading">
                  Current Kubernetes version <span>{currentVersion}</span>
                </p>
              </AlertDialogDescription>

              <div className="change-version-select-container">
                <label
                  htmlFor="versionSelect"
                  className="change-version-select-label"
                >
                  Change to version
                </label>
                <ClusterAlertSelect
                  id="versionSelect"
                  className="change-version-select-input"
                  value={selectedVersion}
                  onChange={e => setSelectedVersion(e.target.value)}
                >
                  {versions.map(version => (
                    <option key={version} value={version}>
                      {version} UnDistro
                    </option>
                  ))}
                </ClusterAlertSelect>
              </div>
            </>
          )}

          <footer>
            <ClusterAlertButton
              disabled={currentVersion === selectedVersion}
              isSecondary
              onClick={() => {
                operationType === OperationType.NONE
                  ? defineOperationType()
                  : handleActionConfirm(selectedVersion)
              }}
            >
              {operationType === OperationType.DOWNGRADE && 'Downgrade'}
              {operationType === OperationType.NONE && 'Next'}
              {operationType === OperationType.UPGRADE && 'Upgrade'}
            </ClusterAlertButton>
            <ClusterAlertButton
              ref={leastDestructiveRef}
              onClick={handleDismiss}
            >
              Cancel
            </ClusterAlertButton>
          </footer>
        </AlertDialog>
      )}
    </>
  )
}

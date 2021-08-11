import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogLabel
} from '@reach/alert-dialog'
import { useRef, useState } from 'react'
import { ReactComponent as CriticalIcon } from '@assets/icons/nodepool-actions/critical.svg'
import { NodepoolActionAlertProps } from '../index'
import { NodepoolAlertButton } from '../common/button'
import { NodepoolAlertInput } from '../common/input'
import './index.scss'

enum NodepoolType {
  WORKER,
  INFRANODE
}

type Props = NodepoolActionAlertProps & {
  isLast?: boolean
  type: NodepoolType
}

export const DeleteNodepoolAlert = ({
  heading,
  isLast,
  isOpen,
  type,
  onActionConfirm: handleActionConfirm,
  onDismiss: handleDismiss
}: Props) => {
  const leastDestructiveRef = useRef(null)
  const [text, setText] = useState('')

  return (
    <>
      {isOpen && (
        <AlertDialog leastDestructiveRef={leastDestructiveRef}>
          <AlertDialogLabel>{heading}</AlertDialogLabel>

          <AlertDialogDescription>
            <CriticalIcon />
            <h3 className="action">
              Delete {type === NodepoolType.INFRANODE ? 'Infranode' : 'Workers'}{' '}
              Nodepool
            </h3>
            <p className="danger-message">
              {isLast
                ? 'At least one workers nodepool is required. Deleting this nodepool will cause a cluster failure.'
                : 'This action cannot be undone.'}
            </p>
          </AlertDialogDescription>

          <div className="form">
            <label htmlFor="deleteText">Type "delete" to confirm:</label>
            <div data-placeholder="delete">
              <NodepoolAlertInput
                id="deleteText"
                type="text"
                value={text}
                onChange={e => setText(e.target.value)}
              />
            </div>
          </div>

          <footer>
            <NodepoolAlertButton
              disabled={text !== 'delete'}
              isSecondary
              onClick={handleActionConfirm}
            >
              Confirm
            </NodepoolAlertButton>
            <NodepoolAlertButton
              ref={leastDestructiveRef}
              onClick={handleDismiss}
            >
              Cancel
            </NodepoolAlertButton>
          </footer>
        </AlertDialog>
      )}
    </>
  )
}

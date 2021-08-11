import { AlertDialogProps } from '@reach/alert-dialog'

import '@reach/dialog/styles.css'
import './index.scss'

export type NodepoolActionAlertProps = Omit<AlertDialogProps, 'children'> & {
  heading: string
  onActionConfirm: () => void
}

export * from './delete'

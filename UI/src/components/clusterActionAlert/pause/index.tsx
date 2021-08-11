import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogLabel,
} from "@reach/alert-dialog";
import { useRef } from "react";
import { ReactComponent as PauseIcon } from "@assets/icons/cluster-actions/default.svg";
import { ClusterActionAlertProps } from "../index";
import { ClusterAlertButton } from "../common/button";

type Props = ClusterActionAlertProps;

export const PauseClusterAlert = ({
  heading,
  isOpen,
  onActionConfirm: handleActionConfirm,
  onDismiss: handleDismiss,
}: Props) => {
  const leastDestructiveRef = useRef(null);

  return (
    <>
      {isOpen && (
        <AlertDialog leastDestructiveRef={leastDestructiveRef}>
          <AlertDialogLabel>{heading}</AlertDialogLabel>

          <AlertDialogDescription>
            <PauseIcon />
            <h3 className="action">Pause</h3>
            <p className="warning-message">
              This cluster will remain active, but UnDistro management will be
              stopped. You can resume it at any time.
            </p>
          </AlertDialogDescription>

          <footer>
            <ClusterAlertButton isSecondary onClick={handleActionConfirm}>
              Pause
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
  );
};

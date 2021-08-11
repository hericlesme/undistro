import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogLabel,
} from "@reach/alert-dialog";
import { useRef } from "react";
import { ReactComponent as ResumeIcon } from "@assets/icons/cluster-actions/resume.svg";
import { ClusterActionAlertProps } from "../index";
import { ClusterAlertButton } from "../common/button";

type Props = ClusterActionAlertProps;

export const ResumeClusterAlert = ({
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
            <ResumeIcon />
            <h3 className="action">Resume</h3>
            <p className="neutral-message">
              Resume UnDistro cluster management.
            </p>
          </AlertDialogDescription>

          <footer>
            <ClusterAlertButton isSecondary onClick={handleActionConfirm}>
              Resume
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

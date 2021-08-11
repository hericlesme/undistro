import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogLabel,
} from "@reach/alert-dialog";
import { useRef, useState } from "react";
import { ReactComponent as CriticalIcon } from "@assets/icons/cluster-actions/critical.svg";
import { ClusterActionAlertProps } from "../index";
import { ClusterAlertButton } from "../common/button";
import { ClusterAlertInput } from "../common/input";
import "./index.scss";

type Props = ClusterActionAlertProps;

export const DeleteClusterAlert = ({
  heading,
  isOpen,
  onActionConfirm: handleActionConfirm,
  onDismiss: handleDismiss,
}: Props) => {
  const leastDestructiveRef = useRef(null);
  const [text, setText] = useState("");

  return (
    <>
      {isOpen && (
        <AlertDialog leastDestructiveRef={leastDestructiveRef}>
          <AlertDialogLabel>{heading}</AlertDialogLabel>

          <AlertDialogDescription>
            <CriticalIcon />
            <h3 className="action">Delete Cluster</h3>
            <p className="danger-message">This action cannot be undone.</p>
          </AlertDialogDescription>

          <div className="form">
            <label htmlFor="deleteText">Type "delete" to confirm:</label>
            <div data-placeholder="delete">
              <ClusterAlertInput
                id="deleteText"
                type="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
              />
            </div>
          </div>

          <footer>
            <ClusterAlertButton
              disabled={text !== "delete"}
              isSecondary
              onClick={handleActionConfirm}
            >
              Confirm
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

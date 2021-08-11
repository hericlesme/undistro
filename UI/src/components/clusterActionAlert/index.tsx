import { AlertDialogProps } from "@reach/alert-dialog";

import "@reach/dialog/styles.css";
import "./index.scss";

export type ClusterActionAlertProps = Omit<AlertDialogProps, "children"> & {
  heading: string;
  onActionConfirm: () => void;
};

export * from "./changeVersion";
export * from "./delete";
export * from "./pause";
export * from "./resume";

import { ComponentPropsWithoutRef, forwardRef } from "react";
import "./index.scss";

type Props = ComponentPropsWithoutRef<"button"> & {
  isSecondary?: boolean;
};

export const ClusterAlertButton = forwardRef<HTMLButtonElement, Props>(
  function ClusterAlertButton(
    { className = "", isSecondary, ...otherProps },
    forwardedRef
  ) {
    return (
      <button
        {...otherProps}
        ref={forwardedRef}
        className={`${className} cluster-alert-button ${
          isSecondary
            ? "cluster-alert-button-secondary"
            : "cluster-alert-button-primary"
        }`}
      />
    );
  }
);

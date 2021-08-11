import { ComponentPropsWithoutRef, forwardRef } from "react";
import "./index.scss";

type Props = ComponentPropsWithoutRef<"input">;

export const ClusterAlertInput = forwardRef<HTMLInputElement, Props>(
  function ClusterAlertInput({ className = "", ...otherProps }, forwardedRef) {
    return (
      <input
        {...otherProps}
        ref={forwardedRef}
        className={`${className} cluster-alert-input`}
      />
    );
  }
);

import { ComponentPropsWithoutRef, forwardRef } from "react";
import "./index.scss";

type Props = ComponentPropsWithoutRef<"select">;

export const ClusterAlertSelect = forwardRef<HTMLSelectElement, Props>(
  function ClusterAlertSelect({ className = "", ...otherProps }, forwardedRef) {
    return (
      <select
        {...otherProps}
        ref={forwardedRef}
        className={`${className} cluster-alert-select`}
      />
    );
  }
);

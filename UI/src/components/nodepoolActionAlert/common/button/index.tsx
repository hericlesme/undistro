import { ComponentPropsWithoutRef, forwardRef } from 'react'
import './index.scss'

type Props = ComponentPropsWithoutRef<'button'> & {
  isSecondary?: boolean
}

export const NodepoolAlertButton = forwardRef<HTMLButtonElement, Props>(
  function NodepoolAlertButton(
    { className = '', isSecondary, ...otherProps },
    forwardedRef
  ) {
    return (
      <button
        {...otherProps}
        ref={forwardedRef}
        className={`${className} nodepool-alert-button ${
          isSecondary
            ? 'nodepool-alert-button-secondary'
            : 'nodepool-alert-button-primary'
        }`}
      />
    )
  }
)

import { ComponentPropsWithoutRef, forwardRef } from 'react'
import './index.scss'

type Props = ComponentPropsWithoutRef<'input'>

export const NodepoolAlertInput = forwardRef<HTMLInputElement, Props>(
  function NodepoolAlertInput({ className = '', ...otherProps }, forwardedRef) {
    return (
      <input
        {...otherProps}
        ref={forwardedRef}
        className={`${className} nodepool-alert-input`}
      />
    )
  }
)

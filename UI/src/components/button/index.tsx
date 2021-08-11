import React, { FC } from 'react'
import Classnames from 'classnames'
import './index.scss'
import { ComponentPropsWithoutRef } from 'react'

type Props = ComponentPropsWithoutRef<'button'> & {
  variant?: string
  size: string
}

const Button: FC<Props> = ({ children, variant = 'primary', size, disabled, onClick, ...otherProps }) => {
  const style = Classnames('button', `button--${variant}`, `button--${size}`)

  return (
    <button {...otherProps} className={style} disabled={disabled} onClick={onClick}>
      {children}
    </button>
  )
}

export default Button

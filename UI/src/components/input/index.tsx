import React, { FC, FormEventHandler } from 'react'
import Classnames from 'classnames'
import './index.scss'

type Props = {
  className?: string
  type?: string
  label?: string
  placeholder?: string
  value?: string | number | undefined
  disabled?: boolean
  validator?: {}
  onChange?: FormEventHandler<HTMLInputElement>
  addButton?: boolean
  handleEvent?: Function
}

const Input: FC<Props> = ({
  className,
  type = 'text',
  label,
  placeholder,
  value,
  disabled,
  validator,
  onChange,
  addButton,
  handleEvent
}) => {
  const style = Classnames(className, 'input', {
    'input--error': validator
  })

  return (
    <div className={style}>
      {label && <label className={disabled ? 'disabled-label' : ''}>{label}</label>}
      <input type={type} value={value} placeholder={placeholder} disabled={disabled} onChange={onChange} />
      {addButton && <i onClick={() => handleEvent?.()} className="icon-plus" />}
    </div>
  )
}

export default Input

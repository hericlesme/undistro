import React, { FC } from 'react'
import Input from '@components/input'
import Select from '@components/select'
import Button from '@components/button'
import { TypeOption } from '../../types/cluster'
import cn from 'classnames'
import './index.scss'

type TypeSlider = {
  title: string
  keyValue: string
  setKeyValue: Function
  value: string
  setValue: Function
  handleAction: Function
  handleClose: Function
  select?: boolean
  taint?: string
  setTaint?: Function
  options?: TypeOption[]
  direction: string
}

const FormSlider: FC<TypeSlider> = ({
  title,
  keyValue,
  setKeyValue,
  value,
  setValue,
  handleAction,
  handleClose,
  select,
  taint,
  setTaint,
  options,
  direction
}) => {
  const style = cn('form-slider-container', `form-slider-container--${direction}`)

  const formKey = (e: React.FormEvent<HTMLInputElement>) => {
    setKeyValue(e.currentTarget.value)
  }

  const formValue = (e: React.FormEvent<HTMLInputElement>) => {
    setValue(e.currentTarget.value)
  }

  const formTaint = (option: TypeOption | null) => {
    setTaint?.(option)
  }

  return (
    <div className={style}>
      <p className="title-slider">{title}</p>

      <div className="form-slider-content">
        <Input type="text" label="key" value={keyValue} onChange={formKey} />
        <Input type="text" label="value" value={value} onChange={formValue} />
        {select && <Select label="taint effect" value={taint} options={options} onChange={formTaint} />}
        <Button size="small" variant="gray" children="add" onClick={() => handleAction()} />
      </div>

      <Button size="small" variant="gray" children="close" onClick={() => handleClose()} />
    </div>
  )
}

export default FormSlider

import React, { FC } from 'react'
import Select from 'react-select'

import './index.scss'

type OptionType = { value: any; label: string }

type Props = {
  className?: string
  label?: string
  options: any
  onChange?: Function
  value?: string
  placeholder?: string
}

const SelectUndistro: FC<Props> = ({ className, label, options, onChange, value, placeholder = 'Select...' }) => {
  const handleChange = (option: any) => {
    if (onChange) {
      onChange(option.value)
    }
  }

  const getCorrectValue = (): OptionType => {
    return options.filter((elm: OptionType) => elm.value === value)[0]
  }

  return (
    <div className={`select ${className}`}>
      <div className="title-section">
        <label>{label}</label>
      </div>

      <Select
        placeholder={placeholder}
        placeholderCSS="react-select-custom-placeholder"
        options={options}
        onChange={handleChange}
        value={value ? getCorrectValue() : (value as any)}
        classNamePrefix="select-container"
      />
    </div>
  )
}

export default SelectUndistro

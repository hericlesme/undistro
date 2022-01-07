import type { InputHTMLAttributes, VFC } from 'react'
import type { Input } from '@/types/forms'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'

type SelectProps = {
  options: (string | number)[]
} & FormActions &
  Input &
  InputHTMLAttributes<HTMLSelectElement>

const Select: VFC<SelectProps> = ({
  label,
  fieldName,
  placeholder,
  options,
  inputSize = 'fit',
  register,
  ...otherProps
}: SelectProps) => {
  let selectProperties = {
    className: classNames(styles.createClusterTextSelect, styles.input100),
    id: fieldName,
    name: fieldName,
    defaultValue: '',
    ...otherProps
  }

  if (register !== undefined) {
    selectProperties = { ...selectProperties, ...register(fieldName) }
  }

  const sizes: { [key: string]: string } = {
    fit: styles.inputFit,
    sm: styles.inputSmall,
    md: styles.inputMedium,
    lg: styles.inputLarge
  }

  const size = sizes[inputSize]

  return (
    <>
      <div className={classNames(styles.inputBlock, size)}>
        <label className={styles.createClusterLabel} htmlFor={fieldName}>
          {label}
        </label>
        <select {...selectProperties}>
          <option value="" disabled hidden>
            {placeholder}
          </option>
          {options &&
            options.map(option => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
        </select>
      </div>
    </>
  )
}

export { Select }

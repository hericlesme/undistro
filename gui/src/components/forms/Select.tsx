import type { VFC } from 'react'
import type { Input } from '@/types/forms'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'

type SelectProps = {
  options: string[]
} & FormActions &
  Input

const Select: VFC<SelectProps> = ({ label, fieldName, placeholder, options, register }: SelectProps) => {
  const selectProperties = {
    className: classNames(styles.createClusterTextSelect, styles.input100),
    id: fieldName,
    name: fieldName,
    defaultValue: '',
    required: true,
    ...register(fieldName, { required: true })
  }

  return (
    <>
      <div className={styles.inputBlock}>
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

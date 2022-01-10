import type { InputHTMLAttributes, VFC } from 'react'
import { useEffect } from 'react'

import classNames from 'classnames'

import { FormActions } from '@/types/utils'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import { Input } from '@/types/forms'

type InputProps = FormActions & Input & InputHTMLAttributes<HTMLInputElement>

const TextInput: VFC<InputProps> = ({
  label,
  placeholder,
  fieldName,
  required,
  type,
  defaultValue,
  inputSize = 'fit',
  register,
  ...otherProps
}: InputProps) => {
  let inputProperties = {
    id: fieldName,
    name: fieldName,
    type: type,
    placeholder: placeholder,
    className: classNames(styles.createClusterTextInput, styles.input100),
    defaultValue: defaultValue,
    ...otherProps
  }

  if (register !== undefined) {
    inputProperties = { ...inputProperties, ...register(fieldName) }
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
      <div className={classNames(classNames(styles.inputBlock, size))}>
        <label className={styles.createClusterLabel} htmlFor={fieldName}>
          {label}
        </label>
        <input {...inputProperties} />
      </div>
    </>
  )
}

export { TextInput }

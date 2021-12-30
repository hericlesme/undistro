import type { VFC } from 'react'

import classNames from 'classnames'

import { FormActions } from '@/types/utils'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'
import { Input } from '@/types/forms'

type TextInputProps = FormActions & Input

const TextInput: VFC<TextInputProps> = ({ label, placeholder, fieldName, required, register }: TextInputProps) => {
  const inputProperties = {
    id: fieldName,
    name: fieldName,
    type: 'text',
    placeholder: placeholder,
    className: classNames(styles.createClusterTextInput, styles.input100),
    ...register(fieldName, { required })
  }

  return (
    <>
      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor={fieldName}>
          {label}
        </label>
        <input {...inputProperties} />
      </div>
    </>
  )
}

export { TextInput }

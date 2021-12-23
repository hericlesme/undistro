import { FieldValues, UseFormRegister, UseFormSetValue } from 'react-hook-form'

export type FormActions = {
  register: UseFormRegister<FieldValues>
  setValue: UseFormSetValue<FieldValues>
}

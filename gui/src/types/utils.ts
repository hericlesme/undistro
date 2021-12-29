import type { Control, FieldValues, UseFormGetValues, UseFormRegister, UseFormSetValue } from 'react-hook-form'

export type FormActions = {
  register: UseFormRegister<FieldValues>
  setValue: UseFormSetValue<FieldValues>
  getValues: UseFormGetValues<FieldValues>
  control: Control<FieldValues, object>
}

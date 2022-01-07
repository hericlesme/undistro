export interface Input {
  label: string
  type?: string
  placeholder: string
  fieldName: string
  required?: boolean
  defaultValue?: string | number | boolean
  inputSize?: string
}

export interface Option {
  value: string
  label: string
}

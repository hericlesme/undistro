export type TypeModal = {
  handleClose: () => void
}

export type TypeSubItem = {
  name: string,
  link?: string,
  handleAction?: () => void
}

export type TypeItem = {
  name: string
  icon: string | any
  subItens: TypeSubItem[]
  url: string
}

export type apiOption = {
  metadata: {
    name: string
  }
  spec: {
    supportedK8SVersions: string[]
  }
}

export type apiResponse = apiOption[]

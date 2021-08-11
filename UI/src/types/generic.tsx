export type TypeModal = {
  handleClose: () => void
}

export type TypeSubItem = {
  name: string
}

export type TypeItem = {
  name: string
  icon: string | any
  subItens: TypeSubItem[]
  url: string
}

export type TypeMenu = {
  itens: TypeItem[]
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

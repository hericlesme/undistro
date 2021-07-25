export type TypeModal = {
  handleClose: () => void
}

export type TypeSubItem = {
  name: string
}

export type TypeItem = {
  name: string
  icon: string
  subItens: TypeSubItem[]
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

export type apiResponse = apiOption[];
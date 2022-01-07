import React, { useState, createContext, useContext } from 'react'
import { ClusterCreation } from '@/components/modals'
import { DeleteCluster } from '@/components/modals/Critical/DeleteCluster'
import { PauseCluster } from '@/components/modals/Critical/PauseCluster'
import { ResumeCluster } from '@/components/modals/Critical/ResumeCluster'

export enum MODAL_TYPES {
  CREATE_CLUSTER = 'CREATE_CLUSTER',
  DELETE_CLUSTER = 'DELETE_CLUSTER',
  PAUSE_CLUSTER = 'PAUSE_CLUSTER',
  RESUME_CLUSTER = 'RESUME_CLUSTER',
  UPDATE_MODAL = 'UPDATE_MODAL'
}

const MODAL_COMPONENTS: Record<string, React.VFC> = {
  [MODAL_TYPES.CREATE_CLUSTER]: ClusterCreation,
  [MODAL_TYPES.DELETE_CLUSTER]: DeleteCluster,
  [MODAL_TYPES.PAUSE_CLUSTER]: PauseCluster,
  [MODAL_TYPES.RESUME_CLUSTER]: ResumeCluster
}

type ModalContext = {
  showModal: (modalType: string, modalProps?: any) => void
  hideModal: () => void
  store: any
}

const initalState: ModalContext = {
  showModal: () => {},
  hideModal: () => {},
  store: {}
}

const ModalContext = createContext(initalState)
export const useModalContext = () => useContext(ModalContext)

type ModalStoreProps = {
  modalType: string
  modalProps: any
}

export const ModalProvider: React.FC<{}> = ({ children }) => {
  const [store, setStore] = useState<ModalStoreProps>()
  const { modalType, modalProps } = store || {}

  const showModal = (modalType: string, modalProps: any = {}) => {
    setStore({
      ...store,
      modalType,
      modalProps
    })
  }

  const hideModal = () => {
    setStore({
      ...store,
      modalType: null,
      modalProps: {}
    })
  }

  const renderComponent = () => {
    const ModalComponent = MODAL_COMPONENTS[modalType]
    if (!modalType || !ModalComponent) {
      return null
    }
    return <ModalComponent id="global-modal" {...modalProps} />
  }

  return (
    <ModalContext.Provider value={{ store, showModal, hideModal }}>
      {renderComponent()}
      {children}
    </ModalContext.Provider>
  )
}

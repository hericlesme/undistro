import { createContext, useContext, ReactNode } from 'react'

import Auth from './api/auth'
import Cluster from './api/cluster'
import Nodepool from './api/nodepool'
import Provider from './api/provider'
import Secret from './api/secret'

type API = {
  Auth: Auth
  Cluster: Cluster
  Nodepool: Nodepool
  Provider: Provider
  Secret: Secret
}

type ServicesContextValue = {
  Api: API
}

type ServicesProviderProps = {
  Api: API
  children: ReactNode
}

const ServicesContext = createContext({} as ServicesContextValue)

export const ServicesProvider = ({ Api, children }: ServicesProviderProps) => {
  return <ServicesContext.Provider value={{ Api }}>{children}</ServicesContext.Provider>
}

export const useServices = () => {
  const { Api } = useContext(ServicesContext)

  return { Api }
}

import { AxiosInstance } from 'axios'
import { createContext, useContext, ReactNode } from 'react'

import Auth from './api/auth'
import Cluster from './api/cluster'
import Nodepool from './api/nodepool'
import Provider from './api/provider'
import Secret from './api/secret'
import Rbac from './api/rbac'

type API = {
  Auth: Auth
  Cluster: Cluster
  Nodepool: Nodepool
  Provider: Provider
  Secret: Secret
  Rbac: Rbac
}

type ServicesContextValue = {
  Api: API
  hasAuthEnabled: boolean
  httpClient: AxiosInstance
}

type ServicesProviderProps = ServicesContextValue & {
  children: ReactNode
}

const ServicesContext = createContext({} as ServicesContextValue)

export const ServicesProvider = ({ Api, children, hasAuthEnabled, httpClient }: ServicesProviderProps) => {
  return <ServicesContext.Provider value={{ Api, hasAuthEnabled, httpClient }}>{children}</ServicesContext.Provider>
}

export const useServices = () => {
  const { Api, hasAuthEnabled, httpClient } = useContext(ServicesContext)

  return { Api, hasAuthEnabled, httpClient }
}

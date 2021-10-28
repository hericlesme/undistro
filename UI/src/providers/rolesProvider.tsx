import { createContext, useState, useContext } from 'react'
import type { ReactNode } from 'react'

type Role = {
  name: string
  type: string
  description: string
}

type RolesContextValue = {
  roles: Role[]
  addRole: (role: Role) => void
}

/*const RolesContext = createContext({} as RolesContextValue)

export const RolesProvider = ({ children } : { children: ReactNode }) => {
  const [roles, setRoles] = useState<Role[]>([])

  const addRole = (role: Role) => {
    setRoles([...roles, role])
  }

  return <RolesContext.Provider value={{ role, addRole }}>{children}</RolesContext.Provider>
}

export const useRoles = () => {
  const { roles, addRole } = useContext(RolesContext)

  return {
    roles,
    addRole,
  }
}*/

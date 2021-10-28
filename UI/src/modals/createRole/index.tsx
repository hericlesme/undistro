import { useState } from 'react'
import Button from '../../components/button'
import Select from '../../components/select'
import Input from '../../components/input'
import store from '../store'
import { useServices } from 'providers/ServicesProvider'
import './index.scss'

type Props = {
  handleClose: () => void
}

function CreateRoleModal({ handleClose }: Props) {
  const body = store.useState((s: any) => s.body)
  const { Api } = useServices()
  const [role, setRole] = useState<string>('')
  const [namespaces, setNamespaces] = useState<string>()
  const [description, setDescription] = useState<string>()
  const [roleName, setRoleName] = useState<string>()
  const [apiGroups, setApiGroups] = useState<string>()
  const [resources, setResources] = useState<string>()
  const [verbs, setVerbs] = useState<string>()
  const [ruleName, setRuleName] = useState<string>()

  const handleAction = () => {
    const data = {
      "apiVersion": "rbac.authorization.k8s.io/v1",
      "kind": role,
      "metadata": {
        "namespace": namespaces,
        "name": roleName,
        "annotations": {
          "description": description,
          "ruleUniqueId": ruleName
        }
      },
      "rules": {
        "apiGroups": apiGroups?.split(",", apiGroups.length),
        "resources": resources?.split(",", resources.length),
        "verbs": verbs?.split(",", verbs.length)
      }
    }

    switch (role) {
      case 'Role':
        return Api.Rbac.createRole(data, namespaces!)
      case 'RoleBiding':
        return Api.Rbac.createRoleBiding(data, namespaces!)
      case 'ClusterRole':
        return Api.Rbac.createClusterRole(data)
      case 'ClusterRoleBiding':
        return Api.Rbac.createClusterRoleBiding(data)
    }
  }

  const formNamespaces = (e: React.FormEvent<HTMLInputElement>) => {
    setNamespaces(e.currentTarget.value)
  }

  const formDescription = (e: React.FormEvent<HTMLInputElement>) => {
    setDescription(e.currentTarget.value)
  }

  const formRoleName = (e: React.FormEvent<HTMLInputElement>) => {
    setRoleName(e.currentTarget.value)
  }

  const formApiGroups = (e: React.FormEvent<HTMLInputElement>) => {
    setApiGroups(e.currentTarget.value)
  }

  const formResources = (e: React.FormEvent<HTMLInputElement>) => {
    setResources(e.currentTarget.value)
  }

  const formVerbs = (e: React.FormEvent<HTMLInputElement>) => {
    setVerbs(e.currentTarget.value)
  }

  const formRuleName = (e: React.FormEvent<HTMLInputElement>) => {
    setRuleName(e.currentTarget.value)
  }


  const options = [
    { label: 'Role', value: 'Role' },
    { label: 'Cluster Role', value: 'ClusterRole' },
    { label: 'Role Biding', value: 'RoleBiding' },
    { label: 'Cluster Role Biding', value: 'ClusterRoleBiding' },
  ]

  return (
    <>
      <header>
        <h3 className="title">
          <span>{body.title}</span>
        </h3>
        <i onClick={handleClose} className="icon-close" />
      </header>
      <section>
        <div className='content'>
          <div className='roles-content'>
            <Input
              placeholder='choose a name for this role'
              label='Role unique ID'
              value={roleName}
              onChange={formRoleName}
            />
            <Input
              placeholder='insert role description'
              label='Description'
              value={description}
              onChange={formDescription}
            />
            <Select
              options={options}
              value={role}
              onChange={setRole}
              className='rbac-select'
              label='Role Type'
            />
            <Input
              placeholder='insert one or more namespaces here'
              label='Namespaces'
              value={namespaces}
              onChange={formNamespaces}
            />

          </div>
          <div className='rules-content'>
            <h3>Rules</h3>
            <Input
            placeholder='insert one or more verbs here'
            label='Rule unique ID'
            value={ruleName}
            onChange={formRuleName}
            />
            <Input
              placeholder='insert one or more groups here'
              label='apiGroups'
              value={apiGroups}
              onChange={formApiGroups}
            />
            <Input
              placeholder='insert one or more resources here'
              label='Resources'
              value={resources}
              onChange={formResources}
            />
            <Input
              placeholder='insert one or more verbs here'
              label='Verbs'
              value={verbs}
              onChange={formVerbs}
            />
          </div>
        </div>
      </section>
      <div className='button-container-rbac'>
        <Button onClick={handleClose} variant='black' size='large' children='Cancel' />
        <Button onClick={handleAction} size='large' children='Save Changes' />
      </div>
    </>
  )
}

export default CreateRoleModal

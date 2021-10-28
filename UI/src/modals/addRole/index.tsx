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

function AddRoleModal({ handleClose }: Props) {
  const body = store.useState((s: any) => s.body)
  const { Api } = useServices()
  const [role, setRole] = useState<string>('')

  const handleAction = () => {
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
            <Select
              placeholder='choose a name for this role'
              label='Role unique ID'
              options={options}
            />
            <Input
              placeholder='insert role description'
              label='Description'
            />
        </div>
      </section>
      <div className='button-container-rbac'>
        <Button onClick={handleClose} variant='black' size='large' children='Cancel' />
        <Button onClick={handleAction} size='large' children='Save Changes' />
      </div>
    </>
  )
}

export default AddRoleModal

import Input from '@components/input'

import './index.scss'

export default function RbacPage() {
  return (
    <div className='rbac-route'>
      <h3>Assign Role</h3>
      <div className='assign-role'>
        <Input placeholder='enter username' />
        <div className='add-role'>
          <i className='icon-add-circle' />
          <p>Add Role</p>
        </div>
      </div>
    </div>
  )
}
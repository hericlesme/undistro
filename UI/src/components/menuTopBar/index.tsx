import { FC } from 'react'
import Cookies from 'js-cookie'
import Logo from '@assets/images/logo.png'
import Modals from 'util/modals'
import { useHistory } from 'react-router'

import './index.scss'

const MenuTopBar: FC = () => {
  const history = useHistory()
  const showModal = () => {
    Modals.show('create-cluster', {
      title: 'Create',
      ndTitle: 'Cluster',
      width: '720',
      height: '600'
    })
  }

  return (
    <div className="menu-top-container">
      <div className="img-container">
        <img alt="undistro-logo" src={Logo} />
      </div>
      <ul style={{ width: '100%' }}>
        <li>
          <p onClick={() => showModal()}>Create</p>
        </li>
        <li>
          <p>Manage</p>
        </li>
        <li>
          <p>Preferences</p>
        </li>
        <li
          style={{ color: '#fff', marginLeft: 'auto' }}
          onClick={() => {
            Cookies.remove('undistro-login')

            history.push('/auth')
          }}
        >
          Logout
        </li>
      </ul>
    </div>
  )
}

export default MenuTopBar

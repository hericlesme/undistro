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
          <p className='disabled'>Modify</p>
        </li>
        <li>
          <p className='disabled'>Manage</p>
        </li>
        <li>
          <p className='disabled'>Preferences</p>
        </li>
        <li>
          <p className='disabled'>About</p>
        </li>
        <li
          className='logout'
          onClick={() => {
            Cookies.remove('undistro-login')

            history.push('/auth')
          }}
        >
          <p>Logout</p>
        </li>
      </ul>
    </div>
  )
}

export default MenuTopBar

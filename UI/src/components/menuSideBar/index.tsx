import { FC } from 'react'
import { useHistory } from 'react-router-dom'
import { TypeMenu } from '../../types/generic'
import './index.scss'

const MenuSideBar: FC<TypeMenu> = ({ itens }) => {
  const history = useHistory()

  return (
    <div className="menu-side-container">
      <ul className="side-itens">
        {itens.map((elm: any) => {
          return (
            <>
              <li className="side-item" onClick={() => history.push(elm.url)}>
                {typeof elm.icon === 'string' ? <i className={elm.icon} /> : elm.icon}
                <p>{elm.name}</p>
              </li>
            </>
          )
        })}
      </ul>
    </div>
  )
}

export default MenuSideBar

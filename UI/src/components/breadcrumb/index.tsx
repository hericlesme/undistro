import { FC } from 'react'
import { useHistory } from 'react-router'
import './index.scss'

type TypeRoutes = {
  name?: any,
  url?: string
}

type TypeBreadCrumb = {
  routes?: TypeRoutes[]
}

const BreadCrumb: FC<TypeBreadCrumb> = ({ routes }) => {
  const history = useHistory()

  const handleRedirect = (url: string) => {
    history.push(url)
  }

  return (
    <div className='bread-crumb'>
      {routes?.map((route, i) => {
        return (
          <li key={route.url} onClick={() => handleRedirect?.(route.url!)}>
            <p>{route.name}</p>
            <i className='icon-arrow-right' />
          </li>
          )
        })}
    </div>
  )
}

export default BreadCrumb

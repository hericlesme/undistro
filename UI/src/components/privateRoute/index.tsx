import { Redirect, Route, RouteProps } from 'react-router-dom'
import Cookies from 'js-cookie'

export const PrivateRoute = ({ children, ...otherProps }: RouteProps) => {
  // #TODO: Fazer requisição para "who am I".
  const isAuthed = Cookies.get('undistro-login')

  return (
    <Route
      {...otherProps}
      render={({ location }) => {
        return !!isAuthed ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: '/auth',
              state: {
                from: location
              }
            }}
          />
        )
      }}
    />
  )
}

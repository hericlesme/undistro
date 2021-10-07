import { Redirect, Route, RouteProps } from 'react-router-dom'

export const PrivateRoute = ({ children, isAuthed, ...otherProps }: RouteProps & { isAuthed: boolean }) => {
  return (
    <Route
      {...otherProps}
      render={({ location }) => {
        return isAuthed ? (
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

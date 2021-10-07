import { Switch, Route, useLocation } from 'react-router-dom'
import AuthRoute from '@routes/auth'
import HomePageRoute from '@routes/home'
import NodepoolsPage from '@routes/nodepool'
import ControlPlanePage from '@routes/controlPlane'
import ClusterRoute from '@routes/cluster'
import WorkerPage from '@routes/worker'
import RbacPage from '@routes/rbacRoles'
import Modals from './modals'
import { ClustersProvider } from 'providers/ClustersProvider'
import { PrivateRoute } from '@components/privateRoute'
import 'styles/app.scss'
import { useEffect, useState } from 'react'
import { useServices } from 'providers/ServicesProvider'
import Cookies from 'js-cookie'
import { useHistory } from 'react-router'

enum AuthStatus {
  AUTHED = 'AUTHED',
  AUTHING = 'AUTHING',
  NOT_AUTHED = 'NOT_AUTHED'
}

export default function App() {
  const [authStatus, setAuthStatus] = useState<AuthStatus>(AuthStatus.AUTHING)
  const { hasAuthEnabled, httpClient } = useServices()
  const location = useLocation()
  const history = useHistory()

  useEffect(() => {
    if (hasAuthEnabled) {
      const whoAmI = async () => {
        const url = `${window.location.protocol}//${window.location.hostname}/_/apis/identity.concierge.pinniped.dev/v1alpha1/whoamirequests`

        const { data } = await httpClient.post(url, {
          apiVersion: 'identity.concierge.pinniped.dev/v1alpha1',
          kind: 'WhoAmIRequest',
          metadata: { creationTimestamp: null },
          spec: {},
          status: { kubernetesUserInfo: { user: { username: '' } } }
        })

        const authStatus = !!data?.status?.kubernetesUserInfo?.user?.username
          ? AuthStatus.AUTHED
          : AuthStatus.NOT_AUTHED

        setAuthStatus(authStatus)
      }

      whoAmI()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  if (authStatus === AuthStatus.NOT_AUTHED) {
    if (window.location.pathname !== '/auth') {
      Cookies.remove('undistro-login')

      history.push('/')
    }
  }

  const isAuthed = authStatus === AuthStatus.AUTHED
  const isAuthing = authStatus === AuthStatus.AUTHING

  return (
    <ClustersProvider>
      <div className="route-container">
        <div className="route-content">
          <Switch>
            {hasAuthEnabled ? (
              !isAuthing && (
                <>
                  <Route exact path="/auth">
                    <AuthRoute isAuthed={isAuthed} isAuthing={isAuthing} />
                  </Route>
                  <PrivateRoute isAuthed={isAuthed} exact path="/">
                    <HomePageRoute />
                  </PrivateRoute>
                  <PrivateRoute isAuthed={isAuthed} exact path="/nodepools">
                    <NodepoolsPage />
                  </PrivateRoute>
                  <PrivateRoute isAuthed={isAuthed} exact path="/controlPlane/:namespace/:clusterName">
                    <ControlPlanePage />
                  </PrivateRoute>
                  <PrivateRoute isAuthed={isAuthed} exact path="/worker/:namespace/:clusterName">
                    <WorkerPage />
                  </PrivateRoute>
                  <PrivateRoute isAuthed={isAuthed} exact path="/cluster/:namespace/:clusterName">
                    <ClusterRoute />
                  </PrivateRoute>
                  <PrivateRoute isAuthed={isAuthed} exact path="/roles">
                    <RbacPage />
                  </PrivateRoute>
                </>
              )
            ) : (
              <>
                <Route exact path="/auth">
                  <AuthRoute isAuthed={isAuthed} isAuthing={isAuthing} />
                </Route>
                <Route component={HomePageRoute} exact path="/" />
                <Route component={NodepoolsPage} exact path="/nodepools" />
                <Route component={ControlPlanePage} exact path="/controlPlane/:namespace/:clusterName" />
                <Route component={WorkerPage} exact path="/worker/:namespace/:clusterName" />
                <Route component={ClusterRoute} exact path="/cluster/:namespace/:clusterName" />
                <Route component={RbacPage} exact path="/roles" />
              </>
            )}
          </Switch>
          <Modals />
        </div>
      </div>
    </ClustersProvider>
  )
}

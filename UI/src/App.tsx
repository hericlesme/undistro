import { Switch, Route } from 'react-router-dom'
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

type AppProps = {
  hasAuthEnabled: boolean
}

export default function App({ hasAuthEnabled }: AppProps) {
  return (
    <ClustersProvider>
      <div className="route-container">
        <div className="route-content">
          <Switch>
            {hasAuthEnabled ? (
              <>
                <Route exact path="/auth" component={AuthRoute} />
                <PrivateRoute exact path="/">
                  <HomePageRoute />
                </PrivateRoute>
                <PrivateRoute exact path="/nodepools">
                  <NodepoolsPage />
                </PrivateRoute>
                <PrivateRoute exact path="/controlPlane/:namespace/:clusterName">
                  <ControlPlanePage />
                </PrivateRoute>
                <PrivateRoute exact path="/worker/:namespace/:clusterName">
                  <WorkerPage />
                </PrivateRoute>
                <PrivateRoute exact path="/cluster/:namespace/:clusterName">
                  <ClusterRoute />
                </PrivateRoute>
                <PrivateRoute exact path="/roles">
                  <RbacPage />
                </PrivateRoute>
              </>
            ) : (
              <>
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

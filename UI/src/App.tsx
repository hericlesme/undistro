import { Switch, Route } from 'react-router-dom'
import HomePageRoute from '@routes/home'
import NodepoolsPage from '@routes/nodepool'
import ControlPlanePage from '@routes/controlPlane'
import ClusterRoute from '@routes/cluster'
import MenuTop from '@components/menuTopBar'
import MenuSideBar from '@components/menuSideBar'
import WorkerPage from '@routes/worker'
import RbacPage from '@routes/rbacRoles'
import Modals from './modals'
import { ClustersProvider } from 'providers/ClustersProvider'
import 'styles/app.scss'

export default function App() {
  return (
    <ClustersProvider>
      <div className="route-container">
        <MenuTop />
        <MenuSideBar />
        <div className="route-content">
          <Switch>
            <Route exact path="/" component={HomePageRoute} />
            <Route exact path="/nodepools" component={NodepoolsPage} />
            <Route exact path="/controlPlane/:namespace/:clusterName" component={ControlPlanePage} />
            <Route exact path="/worker/:namespace/:clusterName" component={WorkerPage} />
            <Route exact path="/cluster/:namespace/:clusterName" component={ClusterRoute} />
            <Route exact path="/roles" component={RbacPage} />
          </Switch>
          <Modals />
        </div>
      </div>
    </ClustersProvider>
  )
}

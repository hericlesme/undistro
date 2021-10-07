import axios, { AxiosInstance, AxiosResponse } from 'axios'
import Cookies from 'js-cookie'
import ReactDOM from 'react-dom'
import { BrowserRouter as Router } from 'react-router-dom'
import { ServicesProvider } from 'providers/ServicesProvider'
import App from './App'
import reportWebVitals from './reportWebVitals'

import Auth from 'providers/ServicesProvider/api/auth'
import Cluster from 'providers/ServicesProvider/api/cluster'
import Nodepool from 'providers/ServicesProvider/api/nodepool'
import Provider from 'providers/ServicesProvider/api/provider'
import Secret from 'providers/ServicesProvider/api/secret'

import '@assets/font-icon/icons.css'
import 'styles/app.scss'

async function createAuthorizationHeader(httpClient: AxiosInstance) {
  const { data } = await httpClient(
    `https://${window.location.hostname}/authcluster?name=management&namespace=undistro-system`,
    {
      withCredentials: true
    }
  )
  const {
    ca,
    credentials: {
      status: { clientCertificateData, clientKeyData }
    },
    endpoint
  } = data

  const authorizationStr = JSON.stringify({ ca, cert: clientCertificateData, endpoint, key: clientKeyData })

  httpClient.defaults.headers.common['Authorization'] = `Bearer ${btoa(authorizationStr)}`
}

;(async function() {
  const BASE_URL = `${window.location.protocol}//${window.location.hostname}/`

  const authConfigUrl = `${BASE_URL}uapi/v1/namespaces/undistro-system/clusters/management/proxy/apis/apps/v1/deployments`
  const { data: apps } = await axios(authConfigUrl)

  const hasAuthEnabled = apps.items.some((app: any) => app.metadata.name.includes('pinniped'))
  const httpClient = axios.create({
    baseURL: hasAuthEnabled
      ? `${BASE_URL}_/`
      : `${BASE_URL}uapi/v1/namespaces/undistro-system/clusters/management/proxy/`,
    timeout: 600000
  })

  httpClient.interceptors.response.use(
    (response: AxiosResponse<any>) => {
      return response
    },
    async (error: any) => {
      if (error.response.data.error === 'securecookie: the value is not valid') {
        await httpClient.get('/logout')
        Cookies.remove('undistro-login')

        window.location.href = '/auth'
      }

      return error
    }
  )

  if (Cookies.get('undistro-login')) {
    await createAuthorizationHeader(httpClient)
  }

  const Api = {
    Auth: new Auth(),
    Cluster: new Cluster(httpClient),
    Nodepool: new Nodepool(httpClient),
    Provider: new Provider(httpClient),
    Secret: new Secret(httpClient)
  }

  ReactDOM.render(
    <Router>
      <ServicesProvider Api={Api} hasAuthEnabled={hasAuthEnabled} httpClient={httpClient}>
        <App />
      </ServicesProvider>
    </Router>,
    document.getElementById('root')
  )

  reportWebVitals()
})()

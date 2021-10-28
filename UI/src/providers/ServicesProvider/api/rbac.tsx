import { AxiosInstance } from 'axios'

class Rbac {
  constructor(public httpClient: AxiosInstance) {}

  async listRoleBindings() {
    const url = `/apis/rbac.authorization.k8s.io/v1/rolebindings`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async listClusterRoleBindings() {
    const url = `/apis/rbac.authorization.k8s.io/v1/clusterrolebindings`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async createRole(data: {}, namespace: string) {
    const url = `/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/roles`
    const res = await this.httpClient.post(url, data)

    return res.data
  }

  async createRoleBiding(data: {}, namespace: string) {
    const url = `/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/rolebindings`
    const res = await this.httpClient.post(url, data)

    return res.data
  }

  async createClusterRole(data: {}) {
    const url = `/apis/rbac.authorization.k8s.io/v1/clusterroles`
    const res = await this.httpClient.post(url, data)

    return res.data
  }

  async createClusterRoleBiding(data: {}) {
    const url = '/apis/rbac.authorization.k8s.io/v1/clusterrolebindings'
    const res = await this.httpClient.post(url, data)

    return res.data
  }
}

export default Rbac

import { AxiosInstance } from 'axios'

class Cluster {
  constructor(public httpClient: AxiosInstance) {}

  async list() {
    const url = `apis/app.undistro.io/v1alpha1/clusters`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async get(namespace: string, clusterName: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/clusters/${clusterName}`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async post(data: {}, namespace: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/clusters`
    const res = await this.httpClient.post(url, data)

    return res.data
  }

  async put(data: {}, namespace: string, name: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/clusters/${name}`
    const res = await this.httpClient.patch(url, data, { headers: { 'Content-Type': 'application/merge-patch+json' } })

    return res.data
  }

  async delete(namespace: string, name: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/clusters/${name}`
    const res = await this.httpClient.delete(url)

    return res.data
  }

  async postPolicies(data: {}, namespace: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/defaultpolicies`
    const res = await this.httpClient.post(url, data)

    return res.data
  }

  async postIdentity(data: {}, namespace: string) {
    const url = `apis/app.undistro.io/v1alpha1/namespaces/${namespace}/identities`
    const res = await this.httpClient.post(url, data)

    return res.data
  }
}

export default Cluster

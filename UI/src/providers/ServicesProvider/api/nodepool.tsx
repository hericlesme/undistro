import { AxiosInstance } from 'axios'

class Node {
  constructor(public httpClient: AxiosInstance) {}

  async list() {
    const url = `apis/cluster.x-k8s.io/v1alpha4/machinepools`
    const res = await this.httpClient.get(encodeURI(url))

    return res.data
  }

  async get(namespace: string, name: string) {
    const url = `apis/cluster.x-k8s.io/v1alpha4/namespaces/${namespace}/machinepools/${name}`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async delete(namespace: string, clusterName: string) {
    const url = `apis/cluster.x-k8s.io/v1alpha4/namespaces/${namespace}/machinepools?labelSelector=cluster.x-k8s.io/cluster-name=${clusterName}`
    const res = await this.httpClient.delete(encodeURI(url))

    return res.data
  }
}

export default Node

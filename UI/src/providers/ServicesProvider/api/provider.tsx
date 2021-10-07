import { AxiosInstance } from 'axios'

class Provider {
  constructor(public httpClient: AxiosInstance) {}

  async list(kind: string) {
    const url = `apis/metadata.undistro.io/v1alpha1/${kind}`
    const res = await this.httpClient.get(url)

    return res.data
  }

  async getUserIp() {
    const url = 'https://api.ipify.org/?format=json'
    const res = await this.httpClient.get(url)

    return res.data
  }
}

export default Provider

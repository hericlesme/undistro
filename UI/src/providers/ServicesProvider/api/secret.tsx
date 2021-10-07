import { AxiosInstance } from 'axios'

class Secret {
  constructor(public httpClient: AxiosInstance) {}

  async list(secretRef: string) {
    const url = `api/v1/namespaces/undistro-system/secrets/${secretRef}`
    const res = await this.httpClient.get(url)

    return res.data
  }
}

export default Secret

import axios, { AxiosInstance } from 'axios'

export default class Auth {
  httpClient: AxiosInstance

  constructor() {
    this.httpClient = axios.create()
  }

  async getProviders() {
    const openIdConfigurationUrl = `https://${window.location.hostname}/auth/.well-known/openid-configuration`
    const { data: openIdConfiguration } = await axios(openIdConfigurationUrl)

    const providersUrl =
      openIdConfiguration['discovery.supervisor.pinniped.dev/v1alpha1']['pinniped_identity_providers_endpoint']
    const { data: identityProviders } = await axios(providersUrl)

    const providers = identityProviders['pinniped_identity_providers'].map(({ name }: { name: string }) => {
      return name.split('-')[1]
    }) as string[]

    return providers
  }
}

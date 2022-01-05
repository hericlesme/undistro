import * as k8s from '@kubernetes/client-node'
import * as request from 'request'
import { RESOURCE_PATH } from '@/helpers/constants'

export function getResourcePath({
  server,
  kind,
  resource
}: {
  server: string
  kind: string
  resource: string
}): string {
  return `${server}/${RESOURCE_PATH[kind]}/${resource}`
}

function getKubeClient() {
  const kc = new k8s.KubeConfig()
  kc.loadFromDefault()
  return kc
}

export function getServerAddress(options: request.Options): string {
  const kc = getKubeClient()
  kc.applyToRequest(options)
  return kc.getCurrentCluster().server
}

export async function getSecret({ serviceAccount, namespace }) {
  const kc = getKubeClient()
  const k8sApi = kc.makeApiClient(k8s.CoreV1Api)

  const secretName = await k8sApi
    .readNamespacedServiceAccount(serviceAccount, namespace)
    .then(res => res.body.secrets[0].name)

  const secret = await k8sApi.readNamespacedSecret(secretName, namespace)
  return Buffer.from(secret.body.data.token, 'base64').toString()
}

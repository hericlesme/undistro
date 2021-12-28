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

export function getServerAddress(options: request.Options): string {
  const kc = new k8s.KubeConfig()
  kc.loadFromDefault() // load from default kubeconfig
  kc.applyToRequest(options)
  return kc.getCurrentCluster().server
}

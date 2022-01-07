import * as k8s from '@kubernetes/client-node'
import { formatDuration, getTimeDiffFromNow } from '@/helpers/time'
import { Cluster } from '@/types/cluster'

export const empty: Cluster = {
  name: '',
  provider: '',
  flavor: '',
  k8sVersion: '',
  clusterGroup: '',
  machines: 0,
  status: '',
  age: ''
}

export function getAge(tm: string, humanize = true, compact = true): string | number {
  const diff = getTimeDiffFromNow(tm)
  if (humanize) {
    return formatDuration(diff, compact)
  }
  return diff
}

export const getStatusFromConditions = (conditions: Array<k8s.V1Condition>): string => {
  if (conditions.some(c => c.message.toLowerCase().includes('reconciliation succeeded'))) {
    return 'Ready'
  } else if (conditions.some(c => c.message.toLowerCase() == 'wait cluster to be provisioned')) {
    return 'Provisioning'
  } else if (conditions.some(c => c.message.toLowerCase() == 'paused')) {
    return 'Paused'
  } else if (conditions.some(c => c.message.toLowerCase() == 'deleting')) {
    return 'Deleting'
  } else if (conditions.some(c => c.message.toLowerCase().includes('error'))) {
    return 'Error'
  } else {
    return 'Unknown'
  }
}

import type { NextApiRequest, NextApiResponse } from 'next'
import type { Cluster } from '@/types/cluster'

import { getSession } from 'next-auth/react'
import { Session } from 'next-auth'

import * as request from 'request'

import { isIdentityEnabled } from '@/helpers/identity'
import { DEFAULT_USER_GROUP } from '@/helpers/constants'
import { getResourcePath, getSecret, getServerAddress } from '@/helpers/server'
import axios from 'axios'
import { KubernetesListObject, KubernetesObject } from '@kubernetes/client-node'

type UnDistroSession = Session & {
  user?: {
    groups?: string[]
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Cluster[]>) {
  const { namespace, cluster } = req.query
  const secret = await getSecret({
    namespace: 'undistro-system',
    serviceAccount: 'undistro-controller-manager'
  })

  const opts: any = {
    headers: {
      Authorization: `Bearer ${secret}`,
      'Content-Type': 'application/merge-patch+json'
    }
  }

  if (isIdentityEnabled()) {
    const session: UnDistroSession = await getSession({ req })
    if (session) {
      opts.headers = {
        ...opts.headers,
        'Impersonate-User': session.user.email,
        'Impersonate-Group': session.user.groups.push(DEFAULT_USER_GROUP)
      }
    }
  }

  const server = getServerAddress(opts)

  const baseUrl = getResourcePath({ server: server, kind: 'app', resource: 'namespaces' })
  const url = `${baseUrl}/${namespace}/clusters/${cluster}`

  request.patch({ url: url, body: JSON.stringify(req.body), ...opts }, (error, response, body) => {
    if (error || response.statusCode !== 201) {
      console.log(error)
      let statusCode = response ? response.statusCode : 500
      //@ts-ignore
      return res.status(statusCode).json({ error: JSON.parse(error || body) })
    }

    if (response) {
      res.json(JSON.parse(body))
      return res.status(201).end()
    }
  })
}

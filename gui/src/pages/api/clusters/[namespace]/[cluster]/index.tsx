import type { NextApiRequest, NextApiResponse } from 'next'
import type { Cluster } from '@/types/cluster'

import { getSession } from 'next-auth/react'
import { Session } from 'next-auth'

import * as request from 'request'

import { clusterDataHandler } from '@/helpers/dataFetching'
import { isIdentityEnabled } from '@/helpers/identity'
import { DEFAULT_USER_GROUP } from '@/helpers/constants'
import { getResourcePath, getServerAddress } from '@/helpers/server'

type UnDistroSession = Session & {
  user?: {
    groups?: string[]
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Cluster[]>) {
  let clusters = []

  const opts = {
    url: req.url
  } as request.Options

  if (isIdentityEnabled()) {
    const session: UnDistroSession = await getSession({ req })
    if (session) {
      opts.headers = {
        'Impersonate-User': session.user.email,
        'Impersonate-Group': session.user.groups.push(DEFAULT_USER_GROUP)
      }
    }
  }

  const { namespace, cluster } = req.query
  const server = getServerAddress(opts)
  const baseUrl = getResourcePath({ server: server, kind: 'app', resource: 'namespaces' })
  const url = `${baseUrl}/${namespace}/clusters/${cluster}`

  request.get(url, opts, (error, response, body) => {
    if (error || response.statusCode !== 200) {
      res.status(response.statusCode).json([])
    }
    if (response) {
      res.json(JSON.parse(body))
      res.status(200).end()
    }
  })
}

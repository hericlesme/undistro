import type { NextApiRequest, NextApiResponse } from 'next'
import type { Cluster } from '@/types/cluster'

import { getSession } from 'next-auth/react'
import { Session } from 'next-auth'

import * as request from 'request'

import { clusterDataHandler } from '@/helpers/dataFetching'
import { isIdentityEnabled } from '@/helpers/identity'
import { DEFAULT_USER_GROUP } from '@/helpers/constants'
import { getResourcePath, getServerAddress } from '@/helpers/server'
import { ErrorMessage } from '@/lib/axios'

type UnDistroSession = Session & {
  user?: {
    groups?: string[]
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Cluster[] | ErrorMessage>) {
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

  const server = getServerAddress(opts)
  const url = getResourcePath({ server: server, kind: 'app', resource: 'clusters' })

  request.get(url, opts, (error, response, body) => {
    if (error || response.statusCode !== 200) {
      return res.status(response.statusCode).json({
        type: 'error',
        status: response.statusCode,
        message: JSON.parse(body).error
      })
    }
    if (response) {
      clusters = clusterDataHandler(JSON.parse(body))
      res.json(clusters)
      return res.status(200).end()
    }
  })
}

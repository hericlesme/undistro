import type { NextApiRequest, NextApiResponse } from 'next'
import type { Cluster } from '@/types/cluster'

import { getSession } from 'next-auth/react'
import { Session } from 'next-auth'

import * as request from 'request'

import { machineTypeDataHandler } from '@/helpers/dataFetching'
import { isIdentityEnabled } from '@/helpers/identity'
import { DEFAULT_USER_GROUP } from '@/helpers/constants'
import { getResourcePath, getServerAddress } from '@/helpers/server'

type UnDistroSession = Session & {
  user?: {
    groups?: string[]
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Cluster[]>) {
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
  const url = getResourcePath({ server: server, kind: 'metadata', resource: 'awsmachines' })

  request.get(url, opts, (error, response, body) => {
    if (error || response.statusCode !== 200) {
      res.status(response.statusCode).send(error)
    }

    if (response) {
      let machines = machineTypeDataHandler(JSON.parse(body).items)
      res.json(machines)
      res.status(200).end()
    }
  })
}

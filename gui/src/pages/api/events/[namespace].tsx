import type { NextApiRequest, NextApiResponse } from 'next'
import type { Cluster } from '@/types/cluster'

import { getSession } from 'next-auth/react'
import { Session } from 'next-auth'

import * as request from 'request'

import fetch from 'node-fetch'

import { isIdentityEnabled } from '@/helpers/identity'
import { DEFAULT_USER_GROUP } from '@/helpers/constants'
import { getSecret, getServerAddress } from '@/helpers/server'

type UnDistroSession = Session & {
  user?: {
    groups?: string[]
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse<Cluster[]>) {
  const secret = await getSecret({
    namespace: 'undistro-system',
    serviceAccount: 'undistro-controller-manager'
  })

  const opts = {
    url: req.url,
    headers: { Authorization: `Bearer ${secret}` }
  } as request.Options

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

  let { namespace } = req.query
  console.log(namespace)

  const url = `${server}/api/v1/namespaces/${namespace}/events?watch=1`

  fetch(url, opts)
    .then(response => response.body)
    .then(response =>
      response
        .on('readable', () => {
          let chunk
          while (null !== (chunk = response.read())) {
            console.log(chunk.toString())
            // pipe to response (or whatever)
          }
        })
        .pipe(res)
    )
    .catch(err => console.log(err))
}

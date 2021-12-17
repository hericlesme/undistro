import { isIdentityEnabled } from '@/helpers/identity'
import { NextFetchEvent, NextResponse } from 'next/server'
import type { NextApiRequest } from 'next'

export async function middleware(req: NextApiRequest, ev: NextFetchEvent) {
  let path = new URL(req.url).pathname

  if (!isIdentityEnabled() && path == '/login') {
    return NextResponse.redirect('/', 302)
  }

  return NextResponse.next()
}

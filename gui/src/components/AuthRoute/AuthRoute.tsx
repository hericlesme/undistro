import { useRouter } from 'next/router'
import { useSession } from 'next-auth/react'
import { NextPage } from 'next'
import { isIdentityEnabled } from '@/helpers/identity'

const CheckAuthRoute: NextPage = ({ children }) => {
  const router = useRouter()
  const publicRoutes = ['/login']

  if (isIdentityEnabled()) {
    const { data: session, status } = useSession()
    if (status !== 'loading') {
      if (!session && !publicRoutes.includes(router.route)) {
        router.push('/login')
        return null
      }
    }
  }

  return <>{children}</>
}

export { CheckAuthRoute }

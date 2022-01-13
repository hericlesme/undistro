import type { AppProps } from 'next/app'

import { SessionProvider } from 'next-auth/react'
import { CheckAuthRoute } from '@/components/AuthRoute/AuthRoute'
import { AppProviders } from '@/contexts'

import '@reach/dialog/styles.css'
import '@reach/combobox/styles.css'
import '@/styles/globals.css'
import { useClusters } from '@/contexts/ClusterContext'
import { useEffect } from 'react'

function UndistroDashBoard({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  const { clusters } = useClusters()

  useEffect(() => {
    console.log(clusters)
  }, [clusters])

  return (
    <AppProviders>
      <SessionProvider session={session}>
        <div className="backgroundDefault">
          <CheckAuthRoute>
            <Component {...pageProps} />
          </CheckAuthRoute>
        </div>
      </SessionProvider>
    </AppProviders>
  )
}

export default UndistroDashBoard

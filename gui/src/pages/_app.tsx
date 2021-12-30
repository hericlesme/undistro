import type { AppProps } from 'next/app'

import { SessionProvider } from 'next-auth/react'
import { CheckAuthRoute } from '@/components/AuthRoute/AuthRoute'
import { AppProviders } from '@/contexts'

import '@reach/dialog/styles.css'
import '@/styles/globals.css'
import '@reach/combobox/styles.css'

function UndistroDashBoard({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  return (
    <AppProviders>
      <SessionProvider session={session}>
        <div className="backgroundDefault responsiveScreenHeight">
          <CheckAuthRoute>
            <Component {...pageProps} />
          </CheckAuthRoute>
        </div>
      </SessionProvider>
    </AppProviders>
  )
}

export default UndistroDashBoard

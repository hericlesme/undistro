import { ClientSafeProvider, getProviders, LiteralUnion, signIn } from 'next-auth/react'
import { BuiltInProviderType } from 'next-auth/providers'
import Image from 'next/image'

import logoUnDistroLogin from '@/public/img/logoUnDistroLogin.svg'
import backedByGetup from '@/public/img/backedByGetup.svg'

import styles from './login.module.css'

type LoginProps = {
  providers: Record<LiteralUnion<BuiltInProviderType, string>, ClientSafeProvider>
}

const Login = ({ providers }: LoginProps) => {
  return (
    <>
      <div className={styles.loginPageContainer}>
        <div className={styles.contentTab}>
          <div className={styles.contentTabProviders}>
            <div className={styles.contentTabLogo}>
              <Image src={logoUnDistroLogin} alt="UnDistro Logo" />
            </div>
            <div className={styles.providersSelection}>
              <div className={styles.selectProviderText}>Sign in with</div>
              <div className={styles.providersContainer}>
                {Object.values(providers).map(provider => (
                  <div
                    key={`provider-${provider.id}`}
                    onClick={() => signIn(provider.id, { redirect: true, callbackUrl: '/' })}
                    className={styles.providersIconContainer}
                  >
                    <div className={styles.providerText}>{provider.name}</div>
                    <div className={styles[`provider${provider.name}`]}></div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <div className={styles.contentTabFooter}>
            <Image src={backedByGetup} alt="Backed by GetUp" />
          </div>
        </div>
      </div>
    </>
  )
}

export const getServerSideProps = async () => {
  const providers = await getProviders()
  return {
    props: { providers }
  }
}

export default Login

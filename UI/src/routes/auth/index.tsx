import axios from 'axios'
import Cookies from 'js-cookie'
import { useEffect, useState } from 'react'

import { ReactComponent as GetUpLogo } from '@assets/auth/getup-logo.svg'
import { ReactComponent as UndistroLogo } from '@assets/auth/undistro-logo.svg'
import { ReactComponent as AwsColored } from '@assets/auth/aws-colored.svg'
import { ReactComponent as AwsGray } from '@assets/auth/aws-gray.svg'
import { ReactComponent as GitHubColored } from '@assets/auth/github-colored.svg'
import { ReactComponent as GitHubGray } from '@assets/auth/github-gray.svg'
import { ReactComponent as GoogleColored } from '@assets/auth/google-colored.svg'
import { ReactComponent as GoogleGray } from '@assets/auth/google-gray.svg'
import { ReactComponent as MicrosoftColored } from '@assets/auth/microsoft-colored.svg'
import { ReactComponent as MicrosoftGray } from '@assets/auth/microsoft-gray.svg'
import { ReactComponent as GitLabColored } from '@assets/auth/gitlab-colored.svg'
import { ReactComponent as GitLabGray } from '@assets/auth/gitlab-gray.svg'
import { useServices } from 'providers/ServicesProvider'

import './index.scss'

const PROVIDERS = {
  aws: (
    <>
      <AwsGray />
      <AwsColored />
      <span>Amazon</span>
    </>
  ),
  github: (
    <>
      <GitHubGray />
      <GitHubColored />
      <span>GitHub</span>
    </>
  ),
  google: (
    <>
      <GoogleGray />
      <GoogleColored />
      <span>Google</span>
    </>
  ),
  microsoft: (
    <>
      <MicrosoftGray />
      <MicrosoftColored />
      <span>Microsoft</span>
    </>
  ),
  gitlab: (
    <>
      <GitLabGray />
      <GitLabColored />
      <span>GitLab</span>
    </>
  )
}

const MAX_AUTH_PROVIDERS_PER_ROW = 4

let isFetchingCert = false

const AuthRoute = () => {
  const { Api } = useServices()
  const [providers, setProviders] = useState<string[]>()

  useEffect(() => {
    ;(async () => {
      const providers = await Api.Auth.getProviders()

      setProviders(providers)
    })()
  }, [Api.Auth])

  return (
    <div className="auth-outer-container">
      <div className="auth-body-container">
        <UndistroLogo />
        <div>
          <p className="auth-instruction">Sign in with:</p>
          <div
            className={`auth-signin-methods-container ${
              Number(providers?.length) > MAX_AUTH_PROVIDERS_PER_ROW ? 'compacted-items' : 'centered-items'
            }`}
          >
            {providers?.map(provider => {
              return (
                <button
                  key={provider}
                  className="auth-signin-method"
                  onClick={async () => {
                    const authWindow = window.open(`https://${window.location.hostname}/login?idp=${provider}`)

                    const timer = setInterval(async () => {
                      if (authWindow?.closed) {
                        if (Cookies.get('undistro-login')) {
                          if (!isFetchingCert) {
                            isFetchingCert = true

                            const certUrl = `https://${window.location.hostname}/authcluster?name=management&namespace=undistro-system`
                            const {
                              data: {
                                ca,
                                credentials: {
                                  status: { clientCertificateData, clientKeyData }
                                }
                              }
                            } = await axios(certUrl, { withCredentials: true })

                            Cookies.set('undistro-ca', ca)
                            Cookies.set('undistro-cert', clientCertificateData)
                            Cookies.set('undistro-key', clientKeyData)

                            window.location.href = '/'
                          }
                        } else {
                          clearInterval(timer)
                        }
                      }
                    }, 1000)
                  }}
                >
                  {PROVIDERS[provider as keyof typeof PROVIDERS]}
                </button>
              )
            })}
          </div>
        </div>
        <div />
      </div>
      <div className="auth-logo-container">
        <GetUpLogo />
      </div>
    </div>
  )
}

export default AuthRoute

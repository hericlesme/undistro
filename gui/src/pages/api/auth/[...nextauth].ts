import NextAuth from 'next-auth'
import GitlabProvider from 'next-auth/providers/gitlab'
import GoogleProvider from 'next-auth/providers/google'
import { Provider } from 'next-auth/providers'

import { GITLAB_DISCOVERY_PATH, OAUTH_PERMISSION_SCOPES } from '@/helpers/constants'

const undistroProviders = {
  gitlab: GitlabProvider({
    clientId: process.env.GITLAB_ID,
    clientSecret: process.env.GITLAB_SECRET,
    wellKnown: GITLAB_DISCOVERY_PATH,
    authorization: { params: { scope: OAUTH_PERMISSION_SCOPES } },
    idToken: true,
    profile(profile) {
      return {
        id: profile.sub,
        email: profile.email,
        groups: profile.groups_direct
      }
    }
  }),
  google: GoogleProvider({
    clientId: process.env.GITHUB_ID,
    clientSecret: process.env.GITHUB_SECRET
  })
} as Record<string, Provider>

export default NextAuth({
  providers: [undistroProviders[process.env.IDENTITY_PROVIDER]],
  secret: process.env.SECRET,
  callbacks: {
    jwt: async ({ token, user }) => {
      user && (token.user = user)
      return token
    },
    session: async ({ session, token }) => {
      session.user = token.user
      return session
    }
  }
})

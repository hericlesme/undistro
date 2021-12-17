/** @type {import('next').NextConfig} */
module.exports = {
  typescript: {
    ignoreDevErrors: true,
    ignoreBuildErrors: true
  },
  env: {
    IDENTITY_ENABLED: process.env.IDENTITY_ENABLED
  }
}

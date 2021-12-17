export function isIdentityEnabled(): boolean {
  return process.env.IDENTITY_ENABLED === 'true'
}

export function getIdentityProvider(): string {
  return process.env.IDENTITY_PROVIDER
}

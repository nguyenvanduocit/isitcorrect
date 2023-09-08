let accessToken: string | undefined
export const getAccessToken = async () => {
  if (accessToken) {
    return accessToken
  }

  const resp = await fetch('https://chat.openai.com/api/auth/session')
  if (resp.status === 403) {
    throw new Error('CLOUDFLARE')
  }
  const data = await resp.json().catch(() => ({}))
  if (!data.accessToken) {
    throw new Error('UNAUTHORIZED')
  }

  accessToken = data.accessToken
  return data.accessToken
}

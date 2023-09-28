import { getAccessToken } from './getAccessToken'

export const authFetch = async (
  input: RequestInfo | URL,
  init?: RequestInit
): Promise<Response> => {
  const accessToken = await getAccessToken()
  init = init ?? ({} as RequestInit)
  init.headers = init.headers ?? ({} as HeadersInit)
  init.headers['Authorization'] = 'Bearer ' + accessToken
  return fetch(input, init)
}

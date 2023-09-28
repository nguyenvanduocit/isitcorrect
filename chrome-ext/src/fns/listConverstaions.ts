import { getAccessToken } from './getAccessToken'

export const listConversations = async () => {
  const accessToken = await getAccessToken()
  const response = await fetch(
    'https://chat.openai.com/backend-api/conversations?offset=0&limit=28&order=updated',
    {
      headers: {
        authorization: 'Bearer ' + accessToken
      },
      body: null,
      method: 'GET',
      mode: 'cors',
      credentials: 'include'
    }
  )

  const data = await response.json()
}

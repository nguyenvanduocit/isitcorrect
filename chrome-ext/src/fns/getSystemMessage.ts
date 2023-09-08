import { getAccessToken } from './getAccessToken'

export const getSystemMessage = async () => {
  const accessToken = await getAccessToken()
  const response = await fetch(
    'https://chat.openai.com/backend-api/user_system_messages',
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      }
    }
  )
  return await response.json()
}

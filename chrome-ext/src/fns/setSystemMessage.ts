import { getAccessToken } from './getAccessToken'
import { getSystemMessage } from './getSystemMessage'

const aboutUserMessage = 'My goal is improving my English knowledge.'
const aboutModelMessage =
  'Your act like an English teacher. Response in Markdown format, do not talk much, do not repeat your self or user input.'

export const setSystemMessage = async () => {
  const systemMessage = await getSystemMessage()

  if (
    systemMessage.about_user_message === aboutUserMessage &&
    systemMessage.about_model_message === aboutModelMessage
  ) {
    return
  }

  const accessToken = await getAccessToken()

  await fetch('https://chat.openai.com/backend-api/user_system_messages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${accessToken}`
    },
    body: JSON.stringify({
      about_user_message: 'My goal is improving my English knowledge.',
      about_model_message:
        'Your act like an English teacher. Response in Markdown format, do not talk much, do not repeat your self or user input.',
      enabled: true
    })
  })
  return
}

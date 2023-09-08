import { getAccessToken } from './getAccessToken'
import { fetchSSE } from './fetchSSE'
import { v4 as uuidv4 } from 'uuid'

export const generateAnswer = async (params: GenerateAnswerParams) => {
  const accessToken = await getAccessToken()

  const modelName = 'text-davinci-002-render-sha'
  await fetchSSE('https://chat.openai.com/backend-api/conversation', {
    method: 'POST',
    signal: params.signal,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${accessToken}`
    },
    body: JSON.stringify({
      action: 'next',
      messages: [
        {
          id: params.messageId ?? uuidv4(),
          role: 'user',
          content: {
            content_type: 'text',
            parts: [params.prompt]
          }
        }
      ],
      conversation_id: params.conversationId,
      model: modelName,
      parent_message_id: params.parentMessageId ?? uuidv4()
    }),
    onMessage(message: string) {
      if (message === '[DONE]') {
        params.onEvent({ type: 'done' })
        return
      }

      let data
      try {
        data = JSON.parse(message)
      } catch (err) {
        console.error(err)
        return
      }
      const text = data.message?.content?.parts?.[0]
      if (text) {
        params.onEvent({
          type: 'answer',
          data: {
            text,
            messageId: data.message.id,
            conversationId: data.conversation_id
          }
        })
      }
    }
  })
  return
}

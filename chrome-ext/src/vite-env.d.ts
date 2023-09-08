/// <reference types="vite/client" />

declare type Target =
  | 'youtubeApp'
  | 'chatGptContent'
  | 'youtubeContent'
  | 'background'

declare type VideoDetails = {
  lengthSeconds: number
  title: string
  viewCount: number
}

declare type ChatGptUser = {
  id: string
  name: string
}

declare type ChatGptSession = {
  accessToken: string
  expires: string
  user: ChatGptUser
}

declare type Answer = {
  text: string
  messageId: string
  conversationId: string
}

declare type GenerateEvent =
  | Event<{
      type: 'answer'
      data: Answer
    }>
  | Event<{
      type: 'done'
    }>

declare type GenerateAnswerParams = {
  prompt: string
  onEvent: (event: GenerateEvent) => void
  signal?: AbortSignal
  messageId?: string
  parentMessageId?: string
  conversationId?: string
}

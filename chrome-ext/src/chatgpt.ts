import ReconnectingWebSocket from 'reconnecting-websocket'
import { generateAnswer } from './fns/generateAnswer'
import { setSystemMessage } from './fns/setSystemMessage'

const options = {
  connectionTimeout: 1000
}
const ws = new ReconnectingWebSocket('ws://localhost:8991/ws', [], options)

ws.addEventListener('open', () => {
  console.log('Connected')
})

ws.addEventListener('close', () => {
  console.log('Disconnected')
})

ws.addEventListener('error', (err) => {
  console.error(err)
})

setSystemMessage().then(() => {
  console.log('System message set')
})

let controller: AbortController
ws.addEventListener('message', (e) => {
  if (controller) {
    controller.abort()
  }
  controller = new AbortController()
  generateAnswer({
    prompt: e.data,
    signal: controller.signal,
    async onEvent(event) {
      if (event.type === 'done') {
        ws.send(':done:')
      } else if (event.type === 'answer') {
        ws.send(event.data.text)
      } else if (event.type === 'error') {
        console.error(event.data)
      }
    }
  }).then(() => {
    console.log('Done')
  })
})

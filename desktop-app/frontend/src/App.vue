<script lang="ts" setup>
import {ClipboardGetText, EventsOn, WindowSetAlwaysOnTop} from "../wailsjs/runtime";
import {IsConnected, SendMessage} from "../wailsjs/go/main/App";
import {onMounted, ref, watchEffect} from "vue";
import {Mdit} from "./services/mdit";
const isLoading = ref(false)
const isExtConnected = ref(false)
const message = ref('')
const question = ref('')
const isMonitorClipboard = ref(true)


onMounted(async () => {
  isExtConnected.value = await IsConnected()
})

EventsOn('message', (data: string) => {
  message.value = Mdit(data)
})

EventsOn('ext-connected', (isConnected: boolean) => {
  isExtConnected.value = isConnected
})

setInterval(async () => {
  if (!isMonitorClipboard.value) return

  const text = await ClipboardGetText()
  if (text && text !== question.value) {
    question.value = text
    await sendMessage()
  }
}, 1000)

const onDrop = async (e: DragEvent) => {
  e.preventDefault()
  const text = e.dataTransfer?.getData('text')
  if (text && text !== question.value) {
    question.value = text
    await sendMessage()
  }
}

const sendMessage = async () => {
  isLoading.value = true
  if (question.value) {
    await SendMessage("check the grammar, rewrite to address any issues, and provide concise explanations. If it's already correct, please attempt to make improvements: "+question.value)
  }
  isLoading.value = false
}
</script>

<template>
  <ElContainer>
    <ElMain v-if="isExtConnected" @drop="onDrop">
      <div>
        <el-switch width="100" v-model="isMonitorClipboard" inactive-text="Monitor clipboard" active-text="Stop monitor" inline-prompt />
        <ElInput v-if="!isMonitorClipboard" v-model="question" @keyup.enter="sendMessage" placeholder="Paste your question here" />
      </div>
      <div v-html="message"></div>
    </ElMain>
    <ElMain v-else>
      <p>What you see here is half of the app. To make it work, you need to open Chrome, install the extension, log in, and keep the ChatGPT tab open.</p>
    </ElMain>
  </ElContainer>
</template>

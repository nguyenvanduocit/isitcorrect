<script lang="ts" setup>
import {ClipboardGetText, EventsOn, WindowSetAlwaysOnTop} from "../wailsjs/runtime";
import {IsConnected, SendMessage} from "../wailsjs/go/main/App";
import {onMounted, ref, watchEffect} from "vue";
import {Mdit} from "./services/mdit";
import {Setting} from "@element-plus/icons-vue";
import OptionForm from "./components/OptionForm.vue";

const isLoading = ref(false)
const isExtConnected = ref(false)
const isShowSetting = ref(false)

const answer = ref('')
const question = ref('')


onMounted(async () => {
  isExtConnected.value = await IsConnected()
})

EventsOn('message', (data: string) => {
  answer.value = Mdit(data)
})

EventsOn('message-end', (data: string) => {
  isLoading.value = false
})

EventsOn('ext-connected', (isConnected: boolean) => {
  isExtConnected.value = isConnected
})

EventsOn('question', (data: string) => {
  question.value = data
})

const sendMessage = async () => {
  isLoading.value = true
  isShowSetting.value = false
  answer.value = ''

  if (question.value) {
    await SendMessage(question.value)
  }
}
</script>

<template>
  <ElContainer>
    <ElHeader :class="$style.inner">
      <ElInput v-if="!isShowSetting" v-loading="isLoading" :disabled="isLoading || !isExtConnected" v-model="question"
               @keyup.enter="sendMessage"
               placeholder="Type and press Enter"/>
      <ElText v-else>Settings</ElText>
      <div>
        <ElButton @click.prevent="isShowSetting = !isShowSetting" :icon="Setting" circle></ElButton>
      </div>

    </ElHeader>
    <template v-if="isShowSetting">
      <ElMain>
        <OptionForm/>
      </ElMain>
    </template>
    <template v-else>
      <ElMain v-if="isExtConnected">
        <ElText :class="$style.text" v-html="answer"></ElText>
        <ElText size="small" v-if="isLoading"> typing...</ElText>
      </ElMain>

      <ElMain v-else>
        <p>To make it function, you should open Chrome, install the extension, log in, and keep the ChatGPT tab
          open.</p>
      </ElMain>
    </template>

  </ElContainer>
</template>

<style module lang="sass">
.inner
  width: 100%
  height: auto
  display: flex
  justify-content: space-between
  align-items: center
  padding: 10px 10px
  background: var(--el-color-black)
  color: var(--el-color-white)
  gap: 10px

.text
  word-break: break-word
</style>

<style lang="sass">
ul, ol
  padding-left: 25px
</style>

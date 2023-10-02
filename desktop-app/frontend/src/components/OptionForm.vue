<script setup lang="ts">
import {onMounted, ref} from "vue";
import {GetRoomID, SetRoomID} from "../../wailsjs/go/main/App";

const roomID = ref("")
onMounted(async () => {
  roomID.value = await GetRoomID()
})

const saveOptions = async () => {
  await SetRoomID(roomID.value)
}

</script>

<template>
<ElForm>
  <ElFormItem label="Hotkey" show-message inline-message="At present, we do not offer custom hotkeys, so this is the only option available.">
    <ElInput readonly value="CMD + Shift + C"></ElInput>
  </ElFormItem>
  <ElFormItem  label="Room ID" show-message inline-message="Room ID">
    <ElInput v-model="roomID"></ElInput>
  </ElFormItem>
  <ElButton @click.prevent="saveOptions">Save</ElButton>
</ElForm>
</template>

<style scoped lang="sass">

</style>

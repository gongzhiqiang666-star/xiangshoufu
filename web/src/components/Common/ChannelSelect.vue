<template>
  <el-select
    v-model="selectedValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    :clearable="clearable"
    :loading="loading"
    class="channel-select"
    @change="handleChange"
  >
    <el-option
      v-for="channel in channels"
      :key="channel.id"
      :label="channel.channel_name"
      :value="channel.id"
    />
  </el-select>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getChannelList } from '@/api/channel'
import type { Channel } from '@/types'

interface Props {
  modelValue: number | undefined
  placeholder?: string
  disabled?: boolean
  size?: 'small' | 'default' | 'large'
  clearable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请选择通道',
  disabled: false,
  size: 'default',
  clearable: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | undefined): void
  (e: 'change', channel: Channel | undefined): void
}>()

// 通道列表（从API加载）
const channels = ref<Channel[]>([])
const loading = ref(false)

const selectedValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

function handleChange(value: number | undefined) {
  const channel = channels.value.find((c) => c.id === value)
  emit('change', channel)
}

// 加载通道列表
async function loadChannels() {
  loading.value = true
  try {
    channels.value = await getChannelList()
  } catch (e) {
    console.error('加载通道列表失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadChannels()
})
</script>

<style lang="scss" scoped>
.channel-select {
  min-width: 150px;
}
</style>

<template>
  <el-select
    v-model="selectedValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    :clearable="clearable"
    class="channel-select"
    @change="handleChange"
  >
    <el-option
      v-for="channel in channels"
      :key="channel.id"
      :label="channel.name"
      :value="channel.id"
    />
  </el-select>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Channel {
  id: number
  name: string
  code: string
}

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

// 通道列表（后续可从API获取）
const channels = ref<Channel[]>([
  { id: 1, name: '恒信通', code: 'HENGXINTONG' },
  { id: 2, name: '拉卡拉', code: 'LAKALA' },
  { id: 3, name: '乐刷', code: 'YEAHKA' },
  { id: 4, name: '随行付', code: 'SUIXINGFU' },
  { id: 5, name: '连连支付', code: 'LIANLIAN' },
  { id: 6, name: '杉德支付', code: 'SANDPAY' },
  { id: 7, name: '富友支付', code: 'FUIOU' },
  { id: 8, name: '汇付天下', code: 'HEEPAY' },
])

const selectedValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

function handleChange(value: number | undefined) {
  const channel = channels.value.find((c) => c.id === value)
  emit('change', channel)
}
</script>

<style lang="scss" scoped>
.channel-select {
  min-width: 150px;
}
</style>

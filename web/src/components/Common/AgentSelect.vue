<template>
  <el-select
    v-model="selectedValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    :clearable="clearable"
    :filterable="filterable"
    :remote="remote"
    :remote-method="handleRemoteSearch"
    :loading="loading"
    class="agent-select"
    @change="handleChange"
  >
    <el-option
      v-for="agent in options"
      :key="agent.id"
      :label="formatLabel(agent)"
      :value="agent.id"
    >
      <div class="agent-option">
        <span class="name">{{ agent.name }}</span>
        <span class="code">{{ agent.agent_code }}</span>
        <el-tag v-if="showLevel" size="small" type="info">{{ agent.level }}级</el-tag>
      </div>
    </el-option>
  </el-select>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { searchAgents } from '@/api/agent'
import type { Agent } from '@/types'

interface Props {
  modelValue: number | undefined
  placeholder?: string
  disabled?: boolean
  size?: 'small' | 'default' | 'large'
  clearable?: boolean
  filterable?: boolean
  remote?: boolean
  showLevel?: boolean
  excludeIds?: number[] // 排除的代理商ID
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请选择代理商',
  disabled: false,
  size: 'default',
  clearable: true,
  filterable: true,
  remote: true,
  showLevel: true,
  excludeIds: () => [],
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | undefined): void
  (e: 'change', agent: Agent | undefined): void
}>()

const loading = ref(false)
const options = ref<Agent[]>([])

const selectedValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

// 格式化标签
function formatLabel(agent: Agent) {
  return `${agent.name} (${agent.agent_code})`
}

// 远程搜索
async function handleRemoteSearch(query: string) {
  if (!query) {
    options.value = []
    return
  }

  loading.value = true
  try {
    const data = await searchAgents(query)
    options.value = data.filter((a) => !props.excludeIds.includes(a.id))
  } catch (error) {
    console.error('Search agents error:', error)
    options.value = []
  } finally {
    loading.value = false
  }
}

// 处理选择变化
function handleChange(value: number | undefined) {
  const agent = options.value.find((a) => a.id === value)
  emit('change', agent)
}

// 初始化时加载当前选中的代理商
watch(
  () => props.modelValue,
  async (val) => {
    if (val && !options.value.find((a) => a.id === val)) {
      // 需要加载选中的代理商信息
      loading.value = true
      try {
        const data = await searchAgents(val.toString())
        if (data.length > 0) {
          options.value = [...options.value, ...data]
        }
      } catch (error) {
        console.error('Load agent error:', error)
      } finally {
        loading.value = false
      }
    }
  },
  { immediate: true }
)
</script>

<style lang="scss" scoped>
.agent-select {
  width: 100%;
}

.agent-option {
  display: flex;
  align-items: center;
  gap: $spacing-sm;

  .name {
    font-weight: 500;
  }

  .code {
    color: $text-secondary;
    font-size: 12px;
  }
}
</style>

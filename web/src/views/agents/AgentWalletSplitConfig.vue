<template>
  <el-dialog
    v-model="visible"
    title="钱包拆分配置"
    width="500px"
    :close-on-click-modal="false"
  >
    <el-alert
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    >
      <template #title>
        <strong>重要提示</strong>
      </template>
      <div>
        <p>1. 开启后，分润钱包和服务费钱包将按通道拆分显示</p>
        <p>2. 拆分开关一旦开启<strong>不可关闭</strong></p>
        <p>3. 下级代理商将自动继承此配置</p>
        <p>4. 奖励钱包不受拆分影响</p>
      </div>
    </el-alert>

    <el-form :model="form" label-width="120px">
      <el-form-item label="代理商">
        <span>{{ agentName }}</span>
      </el-form-item>
      <el-form-item label="当前状态">
        <el-tag :type="currentConfig.split_by_channel ? 'success' : 'info'">
          {{ currentConfig.split_by_channel ? '已开启拆分' : '未拆分（汇总显示）' }}
        </el-tag>
      </el-form-item>
      <el-form-item label="拆分设置" v-if="!currentConfig.split_by_channel">
        <el-switch
          v-model="form.split_by_channel"
          active-text="按通道拆分"
          inactive-text="汇总显示"
        />
      </el-form-item>
      <el-form-item v-if="currentConfig.configured_at">
        <span class="config-info">
          配置时间: {{ formatTime(currentConfig.configured_at) }}
        </span>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button
        type="primary"
        :loading="saving"
        :disabled="currentConfig.split_by_channel"
        @click="handleSave"
      >
        {{ currentConfig.split_by_channel ? '已开启（不可修改）' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSplitConfig, setSplitConfig } from '@/api/walletSplit'
import type { AgentWalletSplitConfig } from '@/types/wallet'

const props = defineProps<{
  modelValue: boolean
  agentId: number
  agentName: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}>()

const visible = ref(false)
const saving = ref(false)
const currentConfig = ref<AgentWalletSplitConfig>({
  agent_id: 0,
  split_by_channel: false,
})
const form = reactive({
  split_by_channel: false,
})

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val
    if (val && props.agentId) {
      loadConfig()
    }
  }
)

watch(visible, (val) => {
  emit('update:modelValue', val)
})

async function loadConfig() {
  try {
    const config = await getSplitConfig(props.agentId)
    currentConfig.value = config
    form.split_by_channel = config.split_by_channel
  } catch (error) {
    console.error('Load split config error:', error)
  }
}

async function handleSave() {
  if (!form.split_by_channel) {
    visible.value = false
    return
  }

  try {
    await ElMessageBox.confirm(
      '确定要开启钱包拆分吗？开启后将无法关闭！',
      '确认操作',
      {
        type: 'warning',
        confirmButtonText: '确定开启',
        cancelButtonText: '取消',
      }
    )

    saving.value = true
    await setSplitConfig(props.agentId, { split_by_channel: true })
    ElMessage.success('设置成功')
    emit('saved')
    visible.value = false
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Save split config error:', error)
    }
  } finally {
    saving.value = false
  }
}

function formatTime(time: string | undefined): string {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
.config-info {
  color: #909399;
  font-size: 12px;
}
</style>

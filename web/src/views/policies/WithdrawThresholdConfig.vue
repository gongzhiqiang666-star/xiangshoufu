<template>
  <el-dialog
    v-model="visible"
    title="提现门槛配置"
    width="700px"
    :close-on-click-modal="false"
  >
    <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px">
      提现门槛是代理商提现时的最低金额限制，单位为元。
    </el-alert>

    <el-form :model="form" label-width="120px">
      <el-divider content-position="left">通用门槛</el-divider>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="分润钱包">
            <el-input-number
              v-model="form.profit_threshold"
              :min="1"
              :precision="0"
              style="width: 100%"
            />
            <span class="unit">元</span>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="服务费钱包">
            <el-input-number
              v-model="form.service_threshold"
              :min="1"
              :precision="0"
              style="width: 100%"
            />
            <span class="unit">元</span>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="奖励钱包">
            <el-input-number
              v-model="form.reward_threshold"
              :min="1"
              :precision="0"
              style="width: 100%"
            />
            <span class="unit">元</span>
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">按通道设置（可选）</el-divider>
      <el-table :data="channelThresholds" border size="small">
        <el-table-column prop="channel_name" label="通道" width="120" />
        <el-table-column label="分润钱包门槛">
          <template #default="{ row }">
            <el-input-number
              v-model="row.profit_threshold"
              :min="0"
              :precision="0"
              size="small"
              placeholder="使用通用"
            />
            <span class="unit">元</span>
          </template>
        </el-table-column>
        <el-table-column label="服务费钱包门槛">
          <template #default="{ row }">
            <el-input-number
              v-model="row.service_threshold"
              :min="0"
              :precision="0"
              size="small"
              placeholder="使用通用"
            />
            <span class="unit">元</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="tip">提示：通道门槛为0时，使用通用门槛</div>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getWithdrawThresholds, batchSetWithdrawThresholds } from '@/api/walletSplit'
import type { PolicyWithdrawThreshold, SetWithdrawThresholdRequest } from '@/types/wallet'

const props = defineProps<{
  modelValue: boolean
  templateId: number
  templateName: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}>()

const visible = ref(false)
const saving = ref(false)

// 通用门槛
const form = reactive({
  profit_threshold: 100,
  service_threshold: 50,
  reward_threshold: 100,
})

// 通道列表（示例数据，实际应从API获取）
const channelThresholds = ref([
  { channel_id: 1, channel_name: '恒信通', profit_threshold: 0, service_threshold: 0 },
  { channel_id: 2, channel_name: '拉卡拉', profit_threshold: 0, service_threshold: 0 },
  { channel_id: 3, channel_name: '乐刷', profit_threshold: 0, service_threshold: 0 },
  { channel_id: 4, channel_name: '随行付', profit_threshold: 0, service_threshold: 0 },
])

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val
    if (val && props.templateId) {
      loadThresholds()
    }
  }
)

watch(visible, (val) => {
  emit('update:modelValue', val)
})

async function loadThresholds() {
  try {
    const thresholds = await getWithdrawThresholds(props.templateId)
    // 解析通用门槛
    thresholds.forEach((t: PolicyWithdrawThreshold) => {
      if (t.channel_id === 0) {
        const amountYuan = t.threshold_amount / 100
        if (t.wallet_type === 1) form.profit_threshold = amountYuan
        else if (t.wallet_type === 2) form.service_threshold = amountYuan
        else if (t.wallet_type === 3) form.reward_threshold = amountYuan
      } else {
        // 按通道门槛
        const channel = channelThresholds.value.find(c => c.channel_id === t.channel_id)
        if (channel) {
          const amountYuan = t.threshold_amount / 100
          if (t.wallet_type === 1) channel.profit_threshold = amountYuan
          else if (t.wallet_type === 2) channel.service_threshold = amountYuan
        }
      }
    })
  } catch (error) {
    console.error('Load thresholds error:', error)
  }
}

async function handleSave() {
  const requests: SetWithdrawThresholdRequest[] = []

  // 通用门槛
  requests.push(
    { wallet_type: 1, channel_id: 0, threshold_amount: form.profit_threshold * 100 },
    { wallet_type: 2, channel_id: 0, threshold_amount: form.service_threshold * 100 },
    { wallet_type: 3, channel_id: 0, threshold_amount: form.reward_threshold * 100 }
  )

  // 通道门槛
  channelThresholds.value.forEach(channel => {
    if (channel.profit_threshold > 0) {
      requests.push({
        wallet_type: 1,
        channel_id: channel.channel_id,
        threshold_amount: channel.profit_threshold * 100,
      })
    }
    if (channel.service_threshold > 0) {
      requests.push({
        wallet_type: 2,
        channel_id: channel.channel_id,
        threshold_amount: channel.service_threshold * 100,
      })
    }
  })

  try {
    saving.value = true
    await batchSetWithdrawThresholds(props.templateId, requests)
    ElMessage.success('保存成功')
    emit('saved')
    visible.value = false
  } catch (error) {
    console.error('Save thresholds error:', error)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.unit {
  margin-left: 8px;
  color: #909399;
}
.tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}
</style>

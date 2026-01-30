<template>
  <div class="withdraw-threshold-view">
    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="fetchData" @reset="fetchData">
      <template #extra>
        <el-button type="primary" :loading="saving" @click="handleSaveGeneral">保存通用门槛</el-button>
      </template>
    </SearchForm>

    <!-- 提示信息 -->
    <el-alert
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    >
      <template #title>
        <strong>提现门槛说明</strong>
      </template>
      <div>
        <p>1. 通用门槛适用于所有通道，可按钱包类型分别设置</p>
        <p>2. 按通道门槛优先级高于通用门槛，设置后该通道使用独立门槛</p>
        <p>3. 门槛为0时将使用通用门槛</p>
      </div>
    </el-alert>

    <!-- 通用门槛配置 -->
    <el-card shadow="never" style="margin-bottom: 16px">
      <template #header>
        <span>通用门槛配置</span>
      </template>
      <el-form :model="generalForm" label-width="120px" inline>
        <el-form-item label="分润钱包">
          <el-input-number
            v-model="generalForm.profit_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 150px"
          />
          <span style="margin-left: 8px">元</span>
        </el-form-item>
        <el-form-item label="服务费钱包">
          <el-input-number
            v-model="generalForm.service_fee_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 150px"
          />
          <span style="margin-left: 8px">元</span>
        </el-form-item>
        <el-form-item label="奖励钱包">
          <el-input-number
            v-model="generalForm.reward_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 150px"
          />
          <span style="margin-left: 8px">元</span>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 按通道门槛配置 -->
    <el-card shadow="never">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>按通道门槛配置</span>
          <el-button type="primary" size="small" :icon="Plus" @click="handleAddChannel">添加通道</el-button>
        </div>
      </template>

      <el-table :data="channelThresholds" style="width: 100%">
        <el-table-column prop="channel_name" label="通道" width="150">
          <template #default="{ row }">
            {{ row.channel_name || `通道ID: ${row.channel_id}` }}
          </template>
        </el-table-column>
        <el-table-column label="分润钱包门槛" width="180">
          <template #default="{ row }">
            <el-input-number
              v-model="row.profit_threshold"
              :min="0"
              :precision="2"
              :step="10"
              size="small"
              style="width: 120px"
            />
            <span style="margin-left: 4px">元</span>
          </template>
        </el-table-column>
        <el-table-column label="服务费钱包门槛" width="180">
          <template #default="{ row }">
            <el-input-number
              v-model="row.service_fee_threshold"
              :min="0"
              :precision="2"
              :step="10"
              size="small"
              style="width: 120px"
            />
            <span style="margin-left: 4px">元</span>
          </template>
        </el-table-column>
        <el-table-column label="奖励钱包门槛" width="180">
          <template #default="{ row }">
            <el-input-number
              v-model="row.reward_threshold"
              :min="0"
              :precision="2"
              :step="10"
              size="small"
              style="width: 120px"
            />
            <span style="margin-left: 4px">元</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleSaveChannel(row)">保存</el-button>
            <el-button type="danger" link size="small" @click="handleDeleteChannel(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="channelThresholds.length === 0" description="暂无按通道门槛配置" />
    </el-card>

    <!-- 添加通道对话框 -->
    <el-dialog v-model="addChannelVisible" title="添加通道门槛" width="500px">
      <el-form :model="newChannelForm" label-width="120px">
        <el-form-item label="选择通道" required>
          <ChannelSelect v-model="newChannelForm.channel_id" style="width: 100%" />
        </el-form-item>
        <el-form-item label="分润钱包门槛">
          <el-input-number
            v-model="newChannelForm.profit_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="服务费钱包门槛">
          <el-input-number
            v-model="newChannelForm.service_fee_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="奖励钱包门槛">
          <el-input-number
            v-model="newChannelForm.reward_threshold"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addChannelVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleConfirmAddChannel">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SearchForm from '@/components/Common/SearchForm.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import {
  getWithdrawThresholds,
  setGeneralThresholds,
  setChannelThresholds,
  deleteChannelThreshold,
} from '@/api/withdrawThreshold'

// 通用门槛表单（元）
const generalForm = reactive({
  profit_threshold: 100,
  service_fee_threshold: 50,
  reward_threshold: 100,
})

// 搜索表单（用于SearchForm组件绑定）
const searchForm = reactive({})

// 按通道门槛列表
interface ChannelThresholdRow {
  channel_id: number
  channel_name?: string
  profit_threshold: number
  service_fee_threshold: number
  reward_threshold: number
}
const channelThresholds = ref<ChannelThresholdRow[]>([])

// 加载状态
const loading = ref(false)
const saving = ref(false)

// 添加通道对话框
const addChannelVisible = ref(false)
const newChannelForm = reactive({
  channel_id: undefined as number | undefined,
  profit_threshold: 0,
  service_fee_threshold: 0,
  reward_threshold: 0,
})

// 分转元
function fenToYuan(fen: number): number {
  return fen / 100
}

// 元转分
function yuanToFen(yuan: number): number {
  return Math.round(yuan * 100)
}

// 获取门槛配置
async function fetchData() {
  loading.value = true
  try {
    const res = await getWithdrawThresholds()

    // 解析通用门槛
    for (const t of res.general_thresholds) {
      if (t.wallet_type === 1) {
        generalForm.profit_threshold = fenToYuan(t.threshold_amount)
      } else if (t.wallet_type === 2) {
        generalForm.service_fee_threshold = fenToYuan(t.threshold_amount)
      } else if (t.wallet_type === 3) {
        generalForm.reward_threshold = fenToYuan(t.threshold_amount)
      }
    }

    // 解析按通道门槛，按channel_id分组
    const channelMap = new Map<number, ChannelThresholdRow>()
    for (const t of res.channel_thresholds) {
      if (!channelMap.has(t.channel_id)) {
        channelMap.set(t.channel_id, {
          channel_id: t.channel_id,
          channel_name: t.channel_name,
          profit_threshold: 0,
          service_fee_threshold: 0,
          reward_threshold: 0,
        })
      }
      const row = channelMap.get(t.channel_id)!
      if (t.wallet_type === 1) {
        row.profit_threshold = fenToYuan(t.threshold_amount)
      } else if (t.wallet_type === 2) {
        row.service_fee_threshold = fenToYuan(t.threshold_amount)
      } else if (t.wallet_type === 3) {
        row.reward_threshold = fenToYuan(t.threshold_amount)
      }
    }
    channelThresholds.value = Array.from(channelMap.values())
  } catch (error) {
    console.error('Fetch thresholds error:', error)
  } finally {
    loading.value = false
  }
}

// 保存通用门槛
async function handleSaveGeneral() {
  saving.value = true
  try {
    await setGeneralThresholds({
      profit_threshold: yuanToFen(generalForm.profit_threshold),
      service_fee_threshold: yuanToFen(generalForm.service_fee_threshold),
      reward_threshold: yuanToFen(generalForm.reward_threshold),
    })
    ElMessage.success('通用门槛保存成功')
  } catch (error) {
    console.error('Save general thresholds error:', error)
  } finally {
    saving.value = false
  }
}

// 打开添加通道对话框
function handleAddChannel() {
  newChannelForm.channel_id = undefined
  newChannelForm.profit_threshold = 0
  newChannelForm.service_fee_threshold = 0
  newChannelForm.reward_threshold = 0
  addChannelVisible.value = true
}

// 确认添加通道门槛
async function handleConfirmAddChannel() {
  if (!newChannelForm.channel_id) {
    ElMessage.warning('请选择通道')
    return
  }

  // 检查是否已存在
  if (channelThresholds.value.some(t => t.channel_id === newChannelForm.channel_id)) {
    ElMessage.warning('该通道门槛已存在，请直接编辑')
    return
  }

  saving.value = true
  try {
    await setChannelThresholds({
      channel_id: newChannelForm.channel_id,
      profit_threshold: yuanToFen(newChannelForm.profit_threshold),
      service_fee_threshold: yuanToFen(newChannelForm.service_fee_threshold),
      reward_threshold: yuanToFen(newChannelForm.reward_threshold),
    })
    ElMessage.success('通道门槛添加成功')
    addChannelVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Add channel threshold error:', error)
  } finally {
    saving.value = false
  }
}

// 保存通道门槛
async function handleSaveChannel(row: ChannelThresholdRow) {
  saving.value = true
  try {
    await setChannelThresholds({
      channel_id: row.channel_id,
      profit_threshold: yuanToFen(row.profit_threshold),
      service_fee_threshold: yuanToFen(row.service_fee_threshold),
      reward_threshold: yuanToFen(row.reward_threshold),
    })
    ElMessage.success('通道门槛保存成功')
  } catch (error) {
    console.error('Save channel threshold error:', error)
  } finally {
    saving.value = false
  }
}

// 删除通道门槛
async function handleDeleteChannel(row: ChannelThresholdRow) {
  try {
    await ElMessageBox.confirm(
      `确定要删除该通道的门槛配置吗？删除后将使用通用门槛。`,
      '确认删除',
      { type: 'warning' }
    )

    await deleteChannelThreshold(row.channel_id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    // 用户取消或请求失败
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.withdraw-threshold-view {
  padding: 0;
}
</style>

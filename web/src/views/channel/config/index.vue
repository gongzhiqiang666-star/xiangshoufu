<template>
  <div class="channel-config-container">
    <!-- 通道选择 -->
    <el-card class="channel-selector" shadow="never">
      <template #header>
        <span>通道配置管理</span>
      </template>
      <el-select
        v-model="selectedChannelId"
        placeholder="请选择通道"
        style="width: 300px"
        @change="handleChannelChange"
      >
        <el-option
          v-for="channel in channelList"
          :key="channel.id"
          :label="channel.channel_name"
          :value="channel.id"
        />
      </el-select>
    </el-card>

    <!-- 配置内容 -->
    <el-card v-if="selectedChannelId" class="config-content" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 费率配置 Tab -->
        <el-tab-pane label="费率配置" name="rate">
          <div class="tab-header">
            <el-button type="primary" @click="showAddRateDialog">
              <el-icon><Plus /></el-icon>
              新增费率类型
            </el-button>
          </div>
          <el-table :data="rateConfigs" border stripe>
            <el-table-column prop="rate_code" label="费率编码" width="120" />
            <el-table-column prop="rate_name" label="费率名称" width="120" />
            <el-table-column label="成本范围" width="180">
              <template #default="{ row }">
                {{ row.min_rate }}% ~ {{ row.max_rate }}%
              </template>
            </el-table-column>
            <el-table-column prop="default_rate" label="默认费率" width="100">
              <template #default="{ row }">
                {{ row.default_rate }}%
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="80" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="editRateConfig(row)">编辑</el-button>
                <el-popconfirm title="确定删除吗？" @confirm="deleteRateConfig(row)">
                  <template #reference>
                    <el-button type="danger" link>删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 押金档位 Tab -->
        <el-tab-pane label="押金档位" name="deposit">
          <el-table :data="depositTiers" border stripe>
            <el-table-column prop="tier_code" label="档位编码" width="120" />
            <el-table-column prop="tier_name" label="档位名称" width="120" />
            <el-table-column label="押金金额" width="120">
              <template #default="{ row }">
                {{ formatMoney(row.deposit_amount) }}元
              </template>
            </el-table-column>
            <el-table-column label="返现上限" width="120">
              <template #default="{ row }">
                {{ formatMoney(row.max_cashback_amount) }}元
              </template>
            </el-table-column>
            <el-table-column label="默认返现" width="120">
              <template #default="{ row }">
                {{ formatMoney(row.default_cashback) }}元
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="editDepositTier(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 流量费返现 Tab -->
        <el-tab-pane label="流量费返现" name="sim">
          <div class="tab-header">
            <el-button type="primary" @click="addSimTier">
              <el-icon><Plus /></el-icon>
              添加档位
            </el-button>
            <el-button type="success" @click="saveSimTiers" :loading="savingSimTiers">
              保存配置
            </el-button>
          </div>
          <el-table :data="simCashbackTiers" border stripe>
            <el-table-column prop="tier_order" label="档位序号" width="100" />
            <el-table-column label="档位名称" width="150">
              <template #default="{ row }">
                <el-input v-model="row.tier_name" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="流量费金额" width="130">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.sim_fee_amount"
                  :min="0"
                  :precision="0"
                  size="small"
                  style="width: 100%"
                />
              </template>
            </el-table-column>
            <el-table-column label="返现上限(分)" width="130">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.max_cashback_amount"
                  :min="0"
                  :precision="0"
                  size="small"
                  style="width: 100%"
                />
              </template>
            </el-table-column>
            <el-table-column label="默认返现(分)" width="130">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.default_cashback"
                  :min="0"
                  :max="row.max_cashback_amount"
                  :precision="0"
                  size="small"
                  style="width: 100%"
                />
              </template>
            </el-table-column>
            <el-table-column label="最后档" width="80">
              <template #default="{ row }">
                <el-switch v-model="row.is_last_tier" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right">
              <template #default="{ $index }">
                <el-button type="danger" link @click="removeSimTier($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="tip-text">
            <el-icon><InfoFilled /></el-icon>
            提示：标记为"最后档"的档位将应用于该序号及以后的所有缴费次数
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 费率配置对话框 -->
    <el-dialog
      v-model="rateDialogVisible"
      :title="rateForm.id ? '编辑费率配置' : '新增费率配置'"
      width="500px"
    >
      <el-form :model="rateForm" :rules="rateRules" ref="rateFormRef" label-width="100px">
        <el-form-item label="费率编码" prop="rate_code">
          <el-input v-model="rateForm.rate_code" :disabled="!!rateForm.id" placeholder="如：CREDIT" />
        </el-form-item>
        <el-form-item label="费率名称" prop="rate_name">
          <el-input v-model="rateForm.rate_name" placeholder="如：贷记卡" />
        </el-form-item>
        <el-form-item label="最低成本" prop="min_rate">
          <el-input v-model="rateForm.min_rate" placeholder="如：0.38">
            <template #suffix>%</template>
          </el-input>
        </el-form-item>
        <el-form-item label="最高限制" prop="max_rate">
          <el-input v-model="rateForm.max_rate" placeholder="如：0.60">
            <template #suffix>%</template>
          </el-input>
        </el-form-item>
        <el-form-item label="默认费率" prop="default_rate">
          <el-input v-model="rateForm.default_rate" placeholder="如：0.55">
            <template #suffix>%</template>
          </el-input>
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="rateForm.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRateForm" :loading="submittingRate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 押金档位对话框 -->
    <el-dialog v-model="depositDialogVisible" title="编辑押金档位" width="500px">
      <el-form :model="depositForm" ref="depositFormRef" label-width="100px">
        <el-form-item label="档位名称">
          <el-input :value="depositForm.tier_name" disabled />
        </el-form-item>
        <el-form-item label="押金金额">
          <el-input :value="formatMoney(depositForm.deposit_amount) + '元'" disabled />
        </el-form-item>
        <el-form-item label="返现上限" prop="max_cashback_amount">
          <el-input-number
            v-model="depositForm.max_cashback_amount"
            :min="0"
            :precision="0"
          />
          <span class="unit">分</span>
        </el-form-item>
        <el-form-item label="默认返现" prop="default_cashback">
          <el-input-number
            v-model="depositForm.default_cashback"
            :min="0"
            :max="depositForm.max_cashback_amount"
            :precision="0"
          />
          <span class="unit">分</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="depositDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDepositForm" :loading="submittingDeposit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, InfoFilled } from '@element-plus/icons-vue'
import {
  getChannelList,
  getChannelRateConfigs,
  createChannelRateConfig,
  updateChannelRateConfig,
  deleteChannelRateConfig as deleteRateConfigApi,
  getChannelDepositTiers,
  updateChannelDepositTier,
  getChannelSimCashbackTiers,
  batchSetChannelSimCashbackTiers,
  type ChannelRateConfig,
  type ChannelDepositTier,
  type ChannelSimCashbackTier,
} from '@/api/channel'
import type { Channel } from '@/types'

// 通道列表
const channelList = ref<Channel[]>([])
const selectedChannelId = ref<number | null>(null)
const activeTab = ref('rate')

// 费率配置
const rateConfigs = ref<ChannelRateConfig[]>([])
const rateDialogVisible = ref(false)
const rateFormRef = ref()
const submittingRate = ref(false)
const rateForm = reactive({
  id: 0,
  rate_code: '',
  rate_name: '',
  min_rate: '',
  max_rate: '',
  default_rate: '',
  sort_order: 0,
})
const rateRules = {
  rate_code: [{ required: true, message: '请输入费率编码', trigger: 'blur' }],
  rate_name: [{ required: true, message: '请输入费率名称', trigger: 'blur' }],
  min_rate: [{ required: true, message: '请输入最低成本', trigger: 'blur' }],
  max_rate: [{ required: true, message: '请输入最高限制', trigger: 'blur' }],
}

// 押金档位
const depositTiers = ref<ChannelDepositTier[]>([])
const depositDialogVisible = ref(false)
const submittingDeposit = ref(false)
const depositForm = reactive({
  id: 0,
  tier_name: '',
  deposit_amount: 0,
  max_cashback_amount: 0,
  default_cashback: 0,
})

// 流量费返现档位
const simCashbackTiers = ref<ChannelSimCashbackTier[]>([])
const savingSimTiers = ref(false)

// 格式化金额（分转元）
const formatMoney = (amount: number) => {
  return (amount / 100).toFixed(2)
}

// 加载通道列表
const loadChannelList = async () => {
  try {
    channelList.value = await getChannelList()
    if (channelList.value.length > 0) {
      selectedChannelId.value = channelList.value[0].id
      await handleChannelChange(selectedChannelId.value)
    }
  } catch (error) {
    ElMessage.error('加载通道列表失败')
  }
}

// 切换通道
const handleChannelChange = async (channelId: number) => {
  await Promise.all([
    loadRateConfigs(channelId),
    loadDepositTiers(channelId),
    loadSimCashbackTiers(channelId),
  ])
}

// 加载费率配置
const loadRateConfigs = async (channelId: number) => {
  try {
    rateConfigs.value = await getChannelRateConfigs(channelId)
  } catch (error) {
    ElMessage.error('加载费率配置失败')
  }
}

// 加载押金档位
const loadDepositTiers = async (channelId: number) => {
  try {
    depositTiers.value = await getChannelDepositTiers(channelId)
  } catch (error) {
    ElMessage.error('加载押金档位失败')
  }
}

// 加载流量费返现档位
const loadSimCashbackTiers = async (channelId: number) => {
  try {
    simCashbackTiers.value = await getChannelSimCashbackTiers(channelId)
  } catch (error) {
    ElMessage.error('加载流量费返现档位失败')
  }
}

// 显示新增费率对话框
const showAddRateDialog = () => {
  Object.assign(rateForm, {
    id: 0,
    rate_code: '',
    rate_name: '',
    min_rate: '',
    max_rate: '',
    default_rate: '',
    sort_order: 0,
  })
  rateDialogVisible.value = true
}

// 编辑费率配置
const editRateConfig = (row: ChannelRateConfig) => {
  Object.assign(rateForm, {
    id: row.id,
    rate_code: row.rate_code,
    rate_name: row.rate_name,
    min_rate: row.min_rate,
    max_rate: row.max_rate,
    default_rate: row.default_rate,
    sort_order: row.sort_order,
  })
  rateDialogVisible.value = true
}

// 提交费率表单
const submitRateForm = async () => {
  try {
    await rateFormRef.value.validate()
    submittingRate.value = true

    if (rateForm.id) {
      await updateChannelRateConfig(selectedChannelId.value!, rateForm.id, {
        rate_name: rateForm.rate_name,
        min_rate: rateForm.min_rate,
        max_rate: rateForm.max_rate,
        default_rate: rateForm.default_rate,
        sort_order: rateForm.sort_order,
      })
      ElMessage.success('更新成功')
    } else {
      await createChannelRateConfig(selectedChannelId.value!, {
        rate_code: rateForm.rate_code,
        rate_name: rateForm.rate_name,
        min_rate: rateForm.min_rate,
        max_rate: rateForm.max_rate,
        default_rate: rateForm.default_rate,
        sort_order: rateForm.sort_order,
      })
      ElMessage.success('创建成功')
    }

    rateDialogVisible.value = false
    await loadRateConfigs(selectedChannelId.value!)
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    submittingRate.value = false
  }
}

// 删除费率配置
const deleteRateConfig = async (row: ChannelRateConfig) => {
  try {
    await deleteRateConfigApi(selectedChannelId.value!, row.id)
    ElMessage.success('删除成功')
    await loadRateConfigs(selectedChannelId.value!)
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// 编辑押金档位
const editDepositTier = (row: ChannelDepositTier) => {
  Object.assign(depositForm, {
    id: row.id,
    tier_name: row.tier_name,
    deposit_amount: row.deposit_amount,
    max_cashback_amount: row.max_cashback_amount,
    default_cashback: row.default_cashback,
  })
  depositDialogVisible.value = true
}

// 提交押金表单
const submitDepositForm = async () => {
  try {
    submittingDeposit.value = true
    await updateChannelDepositTier(selectedChannelId.value!, depositForm.id, {
      max_cashback_amount: depositForm.max_cashback_amount,
      default_cashback: depositForm.default_cashback,
    })
    ElMessage.success('更新成功')
    depositDialogVisible.value = false
    await loadDepositTiers(selectedChannelId.value!)
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    submittingDeposit.value = false
  }
}

// 添加流量费档位
const addSimTier = () => {
  const nextOrder = simCashbackTiers.value.length + 1
  simCashbackTiers.value.push({
    id: 0,
    channel_id: selectedChannelId.value!,
    brand_code: '',
    tier_order: nextOrder,
    tier_name: `第${nextOrder}次缴费`,
    is_last_tier: false,
    max_cashback_amount: 0,
    default_cashback: 0,
    sim_fee_amount: 3600,
    status: 1,
  })
}

// 删除流量费档位
const removeSimTier = (index: number) => {
  simCashbackTiers.value.splice(index, 1)
  // 重新排序
  simCashbackTiers.value.forEach((tier, i) => {
    tier.tier_order = i + 1
  })
}

// 保存流量费档位
const saveSimTiers = async () => {
  try {
    savingSimTiers.value = true
    await batchSetChannelSimCashbackTiers(
      selectedChannelId.value!,
      simCashbackTiers.value.map((tier) => ({
        tier_order: tier.tier_order,
        tier_name: tier.tier_name,
        is_last_tier: tier.is_last_tier,
        max_cashback_amount: tier.max_cashback_amount,
        default_cashback: tier.default_cashback,
        sim_fee_amount: tier.sim_fee_amount,
      }))
    )
    ElMessage.success('保存成功')
    await loadSimCashbackTiers(selectedChannelId.value!)
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    savingSimTiers.value = false
  }
}

onMounted(() => {
  loadChannelList()
})
</script>

<style scoped>
.channel-config-container {
  padding: 20px;
}

.channel-selector {
  margin-bottom: 20px;
}

.config-content {
  min-height: 500px;
}

.tab-header {
  margin-bottom: 16px;
  display: flex;
  gap: 10px;
}

.tip-text {
  margin-top: 16px;
  color: #909399;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.unit {
  margin-left: 8px;
  color: #606266;
}
</style>

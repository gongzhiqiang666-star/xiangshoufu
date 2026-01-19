<template>
  <div class="terminal-list-view">
    <PageHeader title="终端管理" sub-title="终端列表">
      <template #extra>
        <el-button :icon="Upload" @click="handleImport">入库</el-button>
        <el-button :icon="Bottom" @click="handleBatchDispatch">批量下发</el-button>
        <el-button :icon="Top" @click="handleBatchRecall">批量回拨</el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">终端总数</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card activated">
          <div class="stat-value">{{ stats.activated }}</div>
          <div class="stat-label">已激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.distributed }}</div>
          <div class="stat-label">未激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card stock">
          <div class="stat-value">{{ stats.stock }}</div>
          <div class="stat-label">库存</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.yesterday_activated }}</div>
          <div class="stat-label">昨日激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card today">
          <div class="stat-value">{{ stats.today_activated }}</div>
          <div class="stat-label">今日激活</div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6" :lg="3">
        <div class="stat-card">
          <div class="stat-value">{{ stats.month_activated }}</div>
          <div class="stat-label">本月激活</div>
        </div>
      </el-col>
    </el-row>

    <!-- 搜索表单 -->
    <SearchForm v-model="searchForm" @search="handleSearch" @reset="handleReset">
      <el-form-item label="通道">
        <ChannelSelect v-model="searchForm.channel_id" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
          <el-option label="库存" value="stock" />
          <el-option label="已下发" value="distributed" />
          <el-option label="已激活" value="activated" />
          <el-option label="已回拨" value="returned" />
        </el-select>
      </el-form-item>
      <el-form-item label="持有代理商">
        <AgentSelect v-model="searchForm.owner_agent_id" />
      </el-form-item>
      <el-form-item label="SN号">
        <el-input v-model="searchForm.sn" placeholder="请输入SN号" clearable />
      </el-form-item>
    </SearchForm>

    <!-- 表格 -->
    <ProTable
      ref="tableRef"
      :data="tableData"
      :loading="loading"
      :total="total"
      v-model:page="page"
      v-model:page-size="pageSize"
      :selection="true"
      @refresh="fetchData"
      @selection-change="handleSelectionChange"
    >
      <el-table-column prop="sn" label="SN号" width="150" />
      <el-table-column prop="channel_name" label="通道" width="100" />
      <el-table-column prop="owner_agent_name" label="持有代理商" width="120" />
      <el-table-column prop="merchant_name" label="绑定商户" width="100">
        <template #default="{ row }">
          {{ row.merchant_name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="activated_at" label="激活时间" width="170">
        <template #default="{ row }">
          {{ row.activated_at || '-' }}
        </template>
      </el-table-column>

      <template #action="{ row }">
        <el-button type="primary" link @click="handleView(row)">详情</el-button>
        <el-button
          v-if="row.status === 'stock'"
          type="success"
          link
          @click="handleDispatch([row])"
        >
          下发
        </el-button>
        <el-button
          v-if="row.status !== 'stock'"
          type="warning"
          link
          @click="handleRecall([row])"
        >
          回拨
        </el-button>
      </template>
    </ProTable>

    <!-- 下发弹窗 -->
    <el-dialog v-model="dispatchDialogVisible" title="终端下发" width="600px">
      <div class="dispatch-info">
        已选终端: {{ selectedTerminals.length }}台
      </div>
      <el-form :model="dispatchForm" label-width="120px">
        <el-form-item label="下发给代理商" required>
          <AgentSelect v-model="dispatchForm.to_agent_id" style="width: 100%" />
        </el-form-item>
        <el-form-item label="设置货款代扣">
          <el-switch v-model="dispatchForm.enable_deduction" />
        </el-form-item>
        <template v-if="dispatchForm.enable_deduction">
          <el-form-item label="单价">
            <el-input-number v-model="dispatchForm.unit_price" :min="0" />
            <span class="form-tip">元/台</span>
          </el-form-item>
          <el-form-item label="总金额">
            <span class="total-amount">
              ¥{{ (dispatchForm.unit_price * selectedTerminals.length).toFixed(2) }}
            </span>
          </el-form-item>
          <el-form-item label="扣款来源">
            <el-checkbox-group v-model="dispatchForm.wallet_sources">
              <el-checkbox label="profit">分润钱包</el-checkbox>
              <el-checkbox label="service">服务费钱包</el-checkbox>
              <el-checkbox label="reward">奖励钱包</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>
      </el-form>
      <div class="dispatch-warning">
        <el-icon><WarningFilled /></el-icon>
        注意: PC端支持跨级下发，APP端仅支持下发给直属下级
      </div>
      <template #footer>
        <el-button @click="dispatchDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitDispatch">确认下发</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Upload, Bottom, Top, WarningFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/Common/PageHeader.vue'
import SearchForm from '@/components/Common/SearchForm.vue'
import ProTable from '@/components/Common/ProTable.vue'
import ChannelSelect from '@/components/Common/ChannelSelect.vue'
import AgentSelect from '@/components/Common/AgentSelect.vue'
import { getTerminals, getTerminalStats, dispatchTerminals } from '@/api/terminal'
import type { Terminal, TerminalStatus, TerminalStats } from '@/types'

const router = useRouter()

// 统计数据
const stats = ref<TerminalStats>({
  total: 0,
  stock: 0,
  distributed: 0,
  activated: 0,
  returned: 0,
  yesterday_activated: 0,
  today_activated: 0,
  month_activated: 0,
})

// 搜索表单
const searchForm = reactive({
  channel_id: undefined as number | undefined,
  status: undefined as TerminalStatus | undefined,
  owner_agent_id: undefined as number | undefined,
  sn: '',
})

// 表格数据
const tableRef = ref()
const tableData = ref<Terminal[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const selectedTerminals = ref<Terminal[]>([])

// 下发弹窗
const dispatchDialogVisible = ref(false)
const dispatchForm = reactive({
  to_agent_id: undefined as number | undefined,
  enable_deduction: false,
  unit_price: 50,
  wallet_sources: ['profit'] as string[],
})

// 状态配置
function getStatusType(status: TerminalStatus) {
  const map: Record<TerminalStatus, string> = {
    stock: 'info',
    distributed: 'warning',
    activated: 'success',
    returned: '',
  }
  return map[status] || ''
}

function getStatusLabel(status: TerminalStatus) {
  const map: Record<TerminalStatus, string> = {
    stock: '库存',
    distributed: '已下发',
    activated: '已激活',
    returned: '已回拨',
  }
  return map[status] || status
}

// 获取统计数据
async function fetchStats() {
  try {
    stats.value = await getTerminalStats()
  } catch (error) {
    console.error('Fetch terminal stats error:', error)
  }
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getTerminals({
      ...searchForm,
      page: page.value,
      page_size: pageSize.value,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Fetch terminals error:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置
function handleReset() {
  page.value = 1
  fetchData()
}

// 选择变化
function handleSelectionChange(selection: Terminal[]) {
  selectedTerminals.value = selection
}

// 查看详情
function handleView(row: Terminal) {
  router.push(`/terminals/${row.id}`)
}

// 入库
function handleImport() {
  ElMessage.info('入库功能开发中...')
}

// 批量下发
function handleBatchDispatch() {
  if (selectedTerminals.value.length === 0) {
    ElMessage.warning('请先选择要下发的终端')
    return
  }
  handleDispatch(selectedTerminals.value)
}

// 下发
function handleDispatch(terminals: Terminal[]) {
  selectedTerminals.value = terminals
  dispatchForm.to_agent_id = undefined
  dispatchForm.enable_deduction = false
  dispatchDialogVisible.value = true
}

// 提交下发
async function handleSubmitDispatch() {
  if (!dispatchForm.to_agent_id) {
    ElMessage.warning('请选择下发的代理商')
    return
  }

  try {
    await dispatchTerminals({
      terminal_ids: selectedTerminals.value.map((t) => t.id),
      to_agent_id: dispatchForm.to_agent_id,
      cargo_deduction: dispatchForm.enable_deduction
        ? {
            unit_price: dispatchForm.unit_price,
            wallet_sources: dispatchForm.wallet_sources,
          }
        : undefined,
    })
    ElMessage.success('下发成功')
    dispatchDialogVisible.value = false
    tableRef.value?.clearSelection()
    fetchData()
    fetchStats()
  } catch (error) {
    console.error('Dispatch terminals error:', error)
  }
}

// 批量回拨
function handleBatchRecall() {
  if (selectedTerminals.value.length === 0) {
    ElMessage.warning('请先选择要回拨的终端')
    return
  }
  handleRecall(selectedTerminals.value)
}

// 回拨
function handleRecall(terminals: Terminal[]) {
  ElMessage.info('回拨功能开发中...')
}

onMounted(() => {
  fetchStats()
  fetchData()
})
</script>

<style lang="scss" scoped>
.terminal-list-view {
  padding: 0;
}

.stats-row {
  margin-bottom: $spacing-md;

  .el-col {
    margin-bottom: $spacing-md;
  }
}

.stat-card {
  background: $bg-white;
  border-radius: $border-radius-md;
  padding: $spacing-md;
  text-align: center;
  box-shadow: $shadow-sm;

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    color: $text-primary;
  }

  .stat-label {
    font-size: 12px;
    color: $text-secondary;
    margin-top: $spacing-xs;
  }

  &.activated .stat-value {
    color: $success-color;
  }

  &.stock .stat-value {
    color: $info-color;
  }

  &.today .stat-value {
    color: $primary-color;
  }
}

.dispatch-info {
  margin-bottom: $spacing-md;
  padding: $spacing-sm;
  background: $bg-color;
  border-radius: $border-radius-sm;
}

.form-tip {
  margin-left: $spacing-sm;
  color: $text-secondary;
}

.total-amount {
  font-size: 18px;
  font-weight: 600;
  color: $primary-color;
}

.dispatch-warning {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  margin-top: $spacing-md;
  padding: $spacing-sm;
  background: lighten($warning-color, 40%);
  border-radius: $border-radius-sm;
  color: $warning-color;
  font-size: 12px;
}
</style>
